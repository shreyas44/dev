package dev

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/shreyas44/dev/db"
)

var ErrNoDevFile = errors.New("no dev.nix file found")

type DevPath string

func (p *DevPath) DB() *db.DB {
	return db.Load(p.dirPath())
}

func (p *DevPath) config() Config {
	var config Config
	evalOut := bytes.NewBuffer(nil)
	cmd := exec.Command("nix", "eval", "-f", path.Join(string(*p), "dev.nix"), "--json")
	cmd.Stdout = evalOut
	cmd.Run()
	json.Unmarshal(evalOut.Bytes(), &config)

	return config
}

func (p *DevPath) dirPath(elem ...string) string {
	return path.Join(string(*p), ".dev-cli", path.Join(elem...))
}

func (p *DevPath) logFilePath(elem ...string) string {
	return path.Join(p.dirPath(), "logs", path.Join(elem...))
}

func (p *DevPath) Init() {
	dirPath := p.dirPath()
	logsPath := p.dirPath("logs")
	devNixPath := p.dirPath("..", "dev.nix")
	nixPath := p.dirPath("nix")
	profilePath := p.dirPath("nix", "profile")

	os.RemoveAll(nixPath)
	mkDirIfNotExists(dirPath)
	mkDirIfNotExists(logsPath)
	mkDirIfNotExists(nixPath)

	initNixEnv(profilePath, devNixPath)
	if config := p.config(); config.Init != "" {
		runInitScript(config.Init)
	}
}

func (p *DevPath) startService(name string, service Service) {
	logFile := p.logFilePath(name + ".log")
	os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	setupEnv(service.Env)

	script := strings.Trim(strings.Trim(service.Cmd, " "), "\n")

	cmd := exec.Command("dev-daemon", name, p.dirPath(), logFile, script)
	cmd.Start()
}

func (p *DevPath) startChild(child Child) {
	panic("TODO")
}

func (p *DevPath) Start() {
	p.Stop()

	config := p.config()
	s := newSpinner("Starting Services", "Started Services")
	s.start()
	defer s.stop()

	for name, service := range config.Services {
		setupEnv(mergeEnvs(config.Env, service.Env))
		p.startService(name, service)
	}

	for _, child := range config.Children {
		p.startChild(child)
	}
}

func (d *DevPath) Stop() {
	s := newSpinner("Stopping Services", "Stopped Services")
	s.start()
	defer s.stop()

	for _, process := range d.DB().Processes {
		if proc, err := os.FindProcess(process.PID); err == nil {
			proc.Signal(os.Interrupt)
		}
	}
}

func mkDirIfNotExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
}

func setupEnv(env map[string]string) {
	for key, value := range env {
		os.Setenv(key, value)
	}
}

func mergeEnvs(envs ...map[string]string) map[string]string {
	result := make(map[string]string)

	for _, env := range envs {
		for key, value := range env {
			result[key] = value
		}
	}

	return result
}

func initNixEnv(profilePath, devNixPath string) {
	s := newSpinner("Installing Dependencies", "Installed Dependencies")
	s.start()
	defer s.stop()

	cmd := exec.Command("nix-env", "--preserve-installed", "-p", profilePath, "-f", devNixPath, "-iA", "deps")
	cmd.Run()

	os.Setenv("PATH", path.Join(profilePath, "bin")+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func runInitScript(script string) {
	s := newSpinner("Running Init Script", "Init Script Completed")
	s.start()
	defer s.stop()

	cmd := exec.Command("bash", "-c", script)
	cmd.Run()
}

func GetDevNixPath(wd string) (DevPath, error) {
	if _, err := os.Stat(path.Join(wd, "dev.nix")); !os.IsNotExist(err) {
		return DevPath(wd), nil
	}

	if wd == "/" {
		return "", ErrNoDevFile
	}

	return GetDevNixPath(path.Dir(wd))
}
