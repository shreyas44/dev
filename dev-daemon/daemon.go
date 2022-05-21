package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/shreyas44/dev/db"
)

func main() {
	var (
		devNixPath = os.Args[1]
		outputFile = os.Args[2]
		script     = os.Args[3]
		db         = db.Load(devNixPath)
		process, _ = db.ProcessByPID(os.Getpid())
		outFile, _ = os.Create(outputFile)
		done       = make(chan bool)
		sig        = make(chan os.Signal, 1)
		cmd        = exec.Command("bash", "-c", script)
	)

	signal.Notify(sig, os.Interrupt)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = outFile
	cmd.Stderr = outFile
	cmd.Start()

	process.Status = "running"
	db.UpdateProcess(process)

	go func() {
		<-sig
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		process.ExitCode = cmd.ProcessState.ExitCode()
		process.Status = "stopped"
		db.UpdateProcess(process)
		done <- true
	}()

	go func() {
		cmd.Wait()
		done <- true
	}()

	<-done
}
