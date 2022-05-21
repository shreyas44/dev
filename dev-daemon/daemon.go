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
		processName = os.Args[1]
		devDir      = os.Args[2]
		outputFile  = os.Args[3]
		script      = os.Args[4]
		outFile, _  = os.Create(outputFile)
		done        = make(chan bool)
		sig         = make(chan os.Signal, 1)
		cmd         = exec.Command("bash", "-c", script)
		process     = db.Process{
			Name:    processName,
			PID:     os.Getpid(),
			LogFile: outputFile,
			Status:  db.ProcessStatusStarting,
		}
	)

	db.UpdateProcess(devDir, process)

	signal.Notify(sig, os.Interrupt)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = outFile
	cmd.Stderr = outFile
	cmd.Start()

	process.Status = "running"
	db.UpdateProcess(devDir, process)

	go func() {
		<-sig
		// cmd.Process.Kill() as syscall.SIGINT doesn't work for whatever reason
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		process.ExitCode = cmd.ProcessState.ExitCode()
		process.Status = "stopped"
		db.UpdateProcess(devDir, process)
		done <- true
	}()

	go func() {
		cmd.Wait()
		done <- true
	}()

	<-done
}
