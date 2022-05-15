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

	"github.com/google/uuid"
)

var ErrNoDevFile = errors.New("no dev.nix file found")

// use --preserve-installed --from-profile

// root option of -f to specify file name
// get path to `nix-env` before running up

// check active nix profile path
// if it's a dev cli profile
// 		switch to another new profile and install all dependencies of older profile and newer profile
// switch to new profile and install all of its dependencies
// create new nix profile with hash of folder path
// keep track of all nix profiles created and which paths they belong to

// on dev down
// check currently active profile
// if current nix profile doesn't include path
// 		warn nix profile not changed and shutdown processes
// if it's a combination of two other profiles
//		switch to the first profile
// else
// 		switch to the older profile

// profile name = hash (path + name +  old profile is dev cli ? hash of old profile : old profile path)

// start background processes with logs to logs/
// keep track of background processes in db.json
// kill background processes and exit nix shell on exit
// add option to use env mode where nix-env is used instead of nix-shell

// const identifier = "com.dev.cli"

type process struct {
	ID      uuid.UUID `json:"id"`
	PID     int       `json:"pid"`
	Name    string    `json:"name"`
	LogFile string    `json:"logFile"`

	// can be current dir or children
	DevPath string `json:"devPath"`
}

type db struct {
	filePath  string
	Processes map[uuid.UUID]process `json:"processes"`
}

func loadDB(dir string) *db {
	filePath := path.Join(dir, "db.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Create(filePath)
		ioutil.WriteFile(filePath, []byte("{}"), os.ModePerm)
	}

	var db db
	data, _ := ioutil.ReadFile(filePath)
	json.Unmarshal(data, &db)
	db.filePath = filePath

	return &db
}

func (db *db) save() {
	data, _ := json.Marshal(db)
	ioutil.WriteFile(db.filePath, data, os.ModePerm)
}

// func (db *db) add(process *process) {
// 	db.Processes = append(db.Processes, *process)
// 	db.save()
// }

// func (config *Config) start() {
// 	// os.UserConfigDir()
// 	// os.Mkdir(cacheDir+"/com.dev.cli", fs.ModeDir)
// 	config.startServices()
// 	config.startSubServices()
// 	defer config.stopServices()
// 	defer config.stopSubServices()

// 	dir := path.Dir(config.path)
// 	shellCmd := exec.Command("nix-shell", path.Join(dir, config.path), "--run", "$SHELL")
// 	shellCmd.Stdin = os.Stdin
// 	shellCmd.Stdout = os.Stdout
// 	shellCmd.Stderr = os.Stderr

// 	shellCmd.Start()
// 	shellCmd.Wait()
// }

// func getDevFile(dir string) string {
// 	p := path.Join(dir, "dev.nix")
// }

type DevPath string

func (p *DevPath) Init() {
	dirPath := path.Join(string(*p), ".dev-cli")
	devNixPath := path.Join(string(*p), "dev.nix")
	logsPath := path.Join(dirPath, "logs")
	nixPath := path.Join(dirPath, "nix")
	profilePath := path.Join(nixPath, "profile")

	os.RemoveAll(nixPath)
	mkDirIfNotExists(dirPath)
	mkDirIfNotExists(logsPath)
	mkDirIfNotExists(nixPath)

	fmt.Println()
	fmt.Println("Installing Dependencies")
	fmt.Println()

	cmd := exec.Command("nix-env", "-p", profilePath, "-f", devNixPath, "-iA", "deps")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	fmt.Println()
	fmt.Println("Installed Dependencies")
	fmt.Println()

	evalOut := bytes.NewBuffer(nil)
	cmd = exec.Command("nix", "eval", "-f", devNixPath, "--json")
	cmd.Stdout = evalOut
	cmd.Run()

	var config Config
	json.Unmarshal(evalOut.Bytes(), &config)
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

// func (p *DevPath) Deactivate() {
// 	rootDB.removeDevPath(*p)
// 	updateBin()
// }

func mkDirIfNotExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
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
