package dev

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/shreyas44/dev/db"
)

var colors = []color.Attribute{
	color.FgBlue,
	color.FgMagenta,
	color.FgCyan,
	color.FgBlack,
	color.FgGreen,
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
	return &ProcessLogEmitter{process, logCh, colors[rand.Intn(len(colors))]}
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

func (l *Logger) watchService(ch chan string, process db.Process) {
	cmd := exec.Command("tail", "-f", "-n", "+1", process.LogFile)
	em := NewProcessLogEmitter(process, ch)
	cmd.Stdout = em
	cmd.Stderr = em

	cmd.Run()
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
