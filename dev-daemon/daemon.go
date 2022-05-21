package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/shreyas44/dev/db"
)

func main() {
	devNixPath := os.Args[1]
	outputFile := os.Args[2]
	script := os.Args[3]
	db := db.Load(devNixPath)
	process, _ := db.ProcessByPID(os.Getpid())
	outFile, _ := os.Create(outputFile)

	cmd := exec.Command("bash", "-c", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = outFile
	cmd.Stderr = outFile

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	cmd.Start()
	process.Status = "running"
	db.UpdateProcess(process)

	go func() {
		<-sig
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		process.ExitCode = cmd.ProcessState.ExitCode()
		process.Status = "stopped"
		db.UpdateProcess(process)
	}()

	cmd.Wait()
}
