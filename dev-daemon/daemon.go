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

	cmd.Run()
}
