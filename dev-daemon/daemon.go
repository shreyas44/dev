package main

import (
	"os"
	"os/exec"
)

func main() {
	outputFile := os.Args[1]
	script := os.Args[2]
	cmd := exec.Command("bash", "-c", script)

	outFile, _ := os.Create(outputFile)
	cmd.Stdout = outFile
	cmd.Stderr = outFile

	// ch := make(chan os.Signal)
	// signal.Notify(ch, os.Interrupt, os.Kill)
	// defer func() {
	// <-ch
	// cmd.Process.Kill()
	// }()

	cmd.Run()
}
