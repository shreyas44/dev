package dev

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"
)

var ErrNoDevFile = errors.New("no dev.nix file found")

// const identifier = "com.dev.cli"

type Markdowner interface {
	Markdown() string
}

type processes []process

func (p *processes) Markdown() string {
	md := "| Name | PID | Started At |\n"
	md += "| ---- | --- | ---------- |\n"

	for _, process := range *p {
		md += process.Markdown() + "\n"
	}

	return md
}

type process struct {
	PID       int       `json:"pid"`
	Name      string    `json:"name"`
	LogFile   string    `json:"logFile"`
	StartedAt time.Time `json:"startedAt"`

	// can be current dir or children
	DevPath DevPath `json:"devPath"`
}

func (p *process) Markdown() string {
	return fmt.Sprintf("| %s | %d | %s |", p.Name, p.PID, p.StartedAt.Format(time.RFC3339))
}

func (p *process) Stop() {
	if proc, err := os.FindProcess(p.PID); err == nil {
		syscall.Kill(-proc.Pid, syscall.SIGKILL)
	}
}

type db struct {
	filePath  string
	Processes map[string]process `json:"processes"`
}

func loadDB(dir string) *db {
	filePath := path.Join(dir, "db.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Create(filePath)
		ioutil.WriteFile(filePath, []byte("{}"), os.ModePerm)
	}

	db := &db{Processes: make(map[string]process)}
	data, _ := ioutil.ReadFile(filePath)
	json.Unmarshal(data, db)
	db.filePath = filePath

	return db
}

func (db *db) save() {
	data, _ := json.Marshal(db)
	ioutil.WriteFile(db.filePath, data, os.ModePerm)
}

func (db *db) addProcess(process *process) {
	db.Processes[process.Name] = *process
	db.save()
}

func (db *db) removeProcess(process *process) {
	delete(db.Processes, process.Name)
	db.save()
}

type DevPath string

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

	fmt.Println()
	fmt.Println("Installing Dependencies")
	fmt.Println()

	cmd := exec.Command("nix-env", "--preserve-installed", "-p", profilePath, "-f", devNixPath, "-iA", "deps")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	fmt.Println()
	fmt.Println("Installed Dependencies")
	fmt.Println()

	config := p.config()
	if config.Init != "" {
		fmt.Println()
		fmt.Println("Running Init Script")
		fmt.Println()

		cmd := exec.Command("bash", "-c", config.Init)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	}
}

func (p *DevPath) startService(name string, service Service) {
	db := loadDB(p.dirPath())
	if process, ok := db.Processes[name]; ok {
		process.Stop()
	}

	logFile := p.logFilePath(name + ".log")
	os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	setupEnv(service.Env)

	script := strings.Trim(strings.Trim(service.Cmd, " "), "\n")

	cmd := exec.Command("dev-daemon", logFile, script)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Start()

	db.addProcess(&process{
		PID:       cmd.Process.Pid,
		Name:      name,
		LogFile:   logFile,
		StartedAt: time.Now(),
		DevPath:   *p,
	})
}

func (p *DevPath) startChild(child Child) {
	panic("TODO")
}

func (p *DevPath) Process(name string) (process, bool) {
	db := loadDB(p.dirPath())
	process, ok := db.Processes[name]
	return process, ok
}

func (p *DevPath) Processes() processes {
	db := loadDB(p.dirPath())
	processes := make(processes, 0, len(db.Processes))
	for _, process := range db.Processes {
		processes = append(processes, process)
	}

	return processes
}

func (p *DevPath) Start() {
	config := p.config()
	setupEnv(config.Env)

	for name, service := range config.Services {
		p.startService(name, service)
	}

	for _, child := range config.Children {
		p.startChild(child)
	}
}

func (d *DevPath) Stop() {
	db := loadDB(d.dirPath())
	for _, process := range db.Processes {
		fmt.Println("stopping process")
		db.removeProcess(&process)
		process.Stop()
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

func GetDevNixPath(wd string) (DevPath, error) {
	if _, err := os.Stat(path.Join(wd, "dev.nix")); !os.IsNotExist(err) {
		return DevPath(wd), nil
	}

	if wd == "/" {
		return "", ErrNoDevFile
	}

	return GetDevNixPath(path.Dir(wd))
}
