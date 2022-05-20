package dev

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
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

func (e *ProcessLogEmitter) Write(p []byte) (n int, err error) {
	c := color.New(e.color)
	str := strings.Trim(strings.Trim(string(p), " "), "\n")
	e.logCh <- fmt.Sprintf("[%s] %s", c.Sprint(e.process.Name), str)
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
	em := &ProcessLogEmitter{process, ch, colors[rand.Intn(len(colors))]}
	cmd.Stdout = em
	cmd.Stderr = em

	cmd.Start()
}

func (l *Logger) Watch() {
	done := make(chan bool)
	logCh := make(chan string)

	for _, p := range l.Processes {
		go l.watchService(logCh, p)
	}

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt)
		<-sig
		done <- true
	}()

	go func() {
		for str := range logCh {
			fmt.Println(str)
		}
	}()

	<-done
}
