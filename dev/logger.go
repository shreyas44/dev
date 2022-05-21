package dev

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

var colors = []color.Attribute{
	color.FgBlue,
	color.FgMagenta,
	color.FgCyan,
	color.FgBlack,
	color.FgGreen,
}

type ProcessLogEmitter struct {
	process
	logCh chan string
	color color.Attribute
}

func NewProcessLogEmitter(process process, logCh chan string) *ProcessLogEmitter {
	return &ProcessLogEmitter{process, logCh, colors[rand.Intn(len(colors))]}
}

func (e *ProcessLogEmitter) Write(p []byte) (int, error) {
	c := color.New(e.color)
	str := strings.Trim(strings.Trim(string(p), " "), "\n")
	for _, line := range strings.Split(str, "\n") {
		e.logCh <- fmt.Sprintf("[%s] %s", c.Sprint(e.process.Name), line)
	}

	return len(p), nil
}

type Logger struct {
	Processes processes
}

func NewLogger(services ...string) *Logger {
	wd, _ := os.Getwd()
	devPath, _ := GetDevNixPath(wd)
	ps := processes{}
	for _, s := range services {
		process, ok := devPath.Process(s)
		if !ok {
			panic("Service not found")
		}

		ps = append(ps, process)
	}

	if len(services) == 0 {
		ps = devPath.Processes()
	}

	return &Logger{ps}
}

func (l *Logger) watchService(ch chan string, process process) {
	cmd := exec.Command("tail", "-f", process.LogFile)
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
