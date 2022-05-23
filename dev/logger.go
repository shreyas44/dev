package dev

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/shreyas44/dev/db"
)

var currentColor = 0
var colors = []color.Attribute{
	color.FgBlue,
	color.FgMagenta,
	color.FgCyan,
	color.FgBlack,
	color.FgGreen,
}

func getColor() color.Attribute {
	currentColor = (currentColor + 1) % len(colors)
	return colors[currentColor]
}

func trim(str string) string {
	return strings.Trim(strings.Trim(str, " "), "\n")
}

type ProcessLogEmitter struct {
	db.Process
	logCh chan string
	color color.Attribute
}

func NewProcessLogEmitter(process db.Process, logCh chan string) *ProcessLogEmitter {
	return &ProcessLogEmitter{process, logCh, getColor()}
}

func (e *ProcessLogEmitter) Write(p []byte) (int, error) {
	c := color.New(e.color)
	str := trim(string(p))
	for _, line := range strings.Split(str, "\n") {
		e.logCh <- c.Sprintf("[%s]", e.Process.Name) + " " + line
	}

	return len(p), nil
}

type Logger struct {
	Processes []db.Process
}

func NewLogger(dev Dev, services ...string) (*Logger, error) {
	ps := []db.Process{}
	for _, s := range services {
		process, ok := dev.DB().ProcessByName(s)
		if !ok {
			return nil, fmt.Errorf("service %s not found", s)
		}

		ps = append(ps, process)
	}

	if len(services) == 0 {
		ps = dev.DB().ProcessesList()
	}

	return &Logger{ps}, nil
}

func (l *Logger) watchService(logCh chan string, process db.Process) {
	em := NewProcessLogEmitter(process, logCh)
	if process.Status == db.ProcessStatusExited {
		cmd := exec.Command("cat", process.LogFile)
		cmd.Stdout = em
		cmd.Stderr = em
		cmd.Run()
		em.Write([]byte(fmt.Sprintf("exited with code %d", process.ExitCode)))
	} else {
		cmd := exec.Command("tail", "-f", "-n", "+1", process.LogFile)
		cmd.Stdout = em
		cmd.Stderr = em
		cmd.Run()
	}
}

func (l *Logger) Watch() {
	logCh := make(chan string)

	for _, p := range l.Processes {
		go l.watchService(logCh, p)
	}

	go func() {
		for str := range logCh {
			fmt.Println(str)
		}
	}()

	immortalize := make(chan bool)
	<-immortalize
}
