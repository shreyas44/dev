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

func getDBPath(dir string) string {
	return path.Join(dir, "db.json")
}

func Load(dir string) *DB {
	filePath := getDBPath(dir)
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

func (db *DB) save() {
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

func Update(dir string, newDb func(DB) DB) {
	filepath := getDBPath(dir)
	fd, _ := syscall.Open(filepath, syscall.O_RDWR, 0)
	syscall.Flock(fd, syscall.LOCK_EX)

	db := Load(dir)
	new := newDb(*db)
	new.save()

	syscall.Flock(fd, syscall.LOCK_UN)
}

func UpdateProcess(dir string, proceess Process) {
	Update(dir, func(db DB) DB {
		db.Processes[proceess.Name] = proceess
		return db
	})
}
