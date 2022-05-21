package dev

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

type Spinner struct {
	finishedLabel string
	spinner       *spinner.Spinner
}

func newSpinner(runningLabel, finishedLabel string) *Spinner {
	return &Spinner{finishedLabel, spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithSuffix(" "+runningLabel))}
}

func (i *Spinner) start() {
	i.spinner.Start()
}

func (i *Spinner) stop() {
	i.spinner.Stop()
	color.Green("✓ %s", i.finishedLabel)
}
