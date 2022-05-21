package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"syscall"
)

type ProcessStatus string

const (
	ProcessStatusStarting ProcessStatus = "starting"
	ProcessStatusRunning  ProcessStatus = "running"
	ProcessStatusStopped  ProcessStatus = "stopped"
	ProcessStatusExited   ProcessStatus = "exited"
)

type Process struct {
	PID      int           `json:"pid"`
	Name     string        `json:"name"`
	LogFile  string        `json:"logFile"`
	Status   ProcessStatus `json:"processStatus"`
	ExitCode int           `json:"exitCode"`
}

type DB struct {
	filePath  string
	Processes map[string]Process `json:"processes"`
}

func Load(dir string) *DB {
	filePath := path.Join(dir, "db.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Create(filePath)
		ioutil.WriteFile(filePath, []byte("{}"), os.ModePerm)
	}

	db := &DB{Processes: make(map[string]Process)}
	data, _ := ioutil.ReadFile(filePath)
	json.Unmarshal(data, db)
	db.filePath = filePath

	return db
}

func (db *DB) Save() {
	data, _ := json.Marshal(db)
	ioutil.WriteFile(db.filePath, data, os.ModePerm)
}

func (db *DB) ProcessByName(name string) (Process, bool) {
	p, ok := db.Processes[name]
	return p, ok
}

func (db *DB) ProcessByPID(pid int) (Process, bool) {
	for _, p := range db.Processes {
		if p.PID == pid {
			return p, true
		}
	}

	return Process{}, false
}

func (db *DB) ProcessesList() []Process {
	processes := []Process{}
	for _, process := range db.Processes {
		processes = append(processes, process)
	}
	return processes
}

func (db *DB) AddProcess(process Process) {
	if _, ok := db.Processes[process.Name]; ok {
		panic("Process already exists")
	}

	db.Processes[process.Name] = process
	db.Save()
}

func (db *DB) UpdateProcess(process Process) {
	db.Processes[process.Name] = process
	db.Save()
}

func (db *DB) RemoveProcess(process Process) {
	delete(db.Processes, process.Name)
	db.Save()
}

func (p *Process) stop() {
	if proc, err := os.FindProcess(p.PID); err == nil {
		syscall.Kill(-proc.Pid, syscall.SIGKILL)
	}
}
