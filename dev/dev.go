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

var DevFileName = "dev.nix"
var ErrNoDevFile = errors.New("no dev.nix file found")

type Dev struct {
	Path string
}

func (d *Dev) DB() *db.DB {
	return db.Load(d.dirPath())
}

func (d *Dev) config() Config {
	var config Config
	evalOut := bytes.NewBuffer(nil)
	cmd := exec.Command("nix", "eval", "-f", path.Join(d.Path, DevFileName), "--json")
	cmd.Stdout = evalOut
	cmd.Run()
	json.Unmarshal(evalOut.Bytes(), &config)

	return config
}

func (d *Dev) dirPath(elem ...string) string {
	return path.Join(d.Path, ".dev", path.Join(elem...))
}

func (d *Dev) logFilePath(elem ...string) string {
	return path.Join(d.dirPath(), "logs", path.Join(elem...))
}

func (d *Dev) Init() {
	dirPath := d.dirPath()
	logsPath := d.dirPath("logs")
	devNixPath := d.dirPath("..", DevFileName)
	nixPath := d.dirPath("nix")
	profilePath := d.dirPath("nix", "profile")

	os.RemoveAll(nixPath)
	mkDirIfNotExists(dirPath)
	mkDirIfNotExists(logsPath)
	mkDirIfNotExists(nixPath)

	d.initNixEnv(profilePath, devNixPath)
	if config := d.config(); config.Init != "" {
		d.runInitScript(config.Init)
	}
}

func (d *Dev) startService(name string, service Service) {
	logFile := d.logFilePath(name + ".log")
	os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	setupEnv(service.Env)

	script := strings.Trim(strings.Trim(service.Cmd, " "), "\n")

	cmd := exec.Command("dev-daemon", name, d.dirPath(), logFile, script)
	cmd.Start()
}

func (d *Dev) startChild(child Child) {
	panic("TODO")
}

func (d *Dev) Start() {
	config := d.config()
	s := newSpinner("Starting Services", "Started Services")
	s.start()
	defer s.stop()

	for name, service := range config.Services {
		setupEnv(mergeEnvs(config.Env, service.Env))
		d.startService(name, service)
	}

	for _, child := range config.Children {
		d.startChild(child)
	}
}

func (d *Dev) Stop() {
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

func (d *Dev) initNixEnv(profilePath, devNixPath string) {
	s := newSpinner("Installing Dependencies", "Installed Dependencies")
	s.start()
	defer s.stop()

	outFile, _ := os.Create(d.dirPath("nix_init.log"))
	cmd := exec.Command("nix-env", "--preserve-installed", "-p", profilePath, "-f", devNixPath, "-iA", "deps")
	cmd.Stdout = outFile
	cmd.Stderr = outFile
	cmd.Run()

	os.Setenv("PATH", path.Join(profilePath, "bin")+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func (d *Dev) runInitScript(script string) {
	s := newSpinner("Running Init Script", "Init Script Completed")
	s.start()
	defer s.stop()

	outFile, _ := os.Create(d.dirPath("init_script.log"))
	cmd := exec.Command("bash", "-c", script)
	cmd.Stdout = outFile
	cmd.Stderr = outFile
	cmd.Run()
}

func Get(wd string) (Dev, error) {
	if _, err := os.Stat(path.Join(wd, DevFileName)); !os.IsNotExist(err) {
		return Dev{wd}, nil
	}

	if wd == "/" {
		return Dev{}, ErrNoDevFile
	}

	return Get(path.Dir(wd))
}
