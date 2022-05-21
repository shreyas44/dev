package dev

import (
	"time"

	bspinner "github.com/briandowns/spinner"
	"github.com/fatih/color"
)

type spinner struct {
	finishedLabel string
	spinner       *bspinner.Spinner
}

func newSpinner(runningLabel, finishedLabel string) *spinner {
	return &spinner{finishedLabel, bspinner.New(bspinner.CharSets[11], 100*time.Millisecond, bspinner.WithSuffix(" "+runningLabel))}
}

func (i *spinner) start() {
	i.spinner.Start()
}

func (i *spinner) stop() {
	i.spinner.Stop()
	color.Green("✓ %s", i.finishedLabel)
}
