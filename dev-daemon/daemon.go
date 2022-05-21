package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	database "github.com/shreyas44/dev/db"
)

func main() {
	var (
		processName = os.Args[1]
		devNixPath  = os.Args[2]
		outputFile  = os.Args[3]
		script      = os.Args[4]
		db          = database.Load(devNixPath)
		outFile, _  = os.Create(outputFile)
		done        = make(chan bool)
		sig         = make(chan os.Signal, 1)
		cmd         = exec.Command("bash", "-c", script)
		process     = database.Process{
			Name:    processName,
			PID:     os.Getpid(),
			LogFile: outputFile,
			Status:  database.ProcessStatusStarting,
		}
	)

	// we add process here to avoid conflicting writes
	db.UpdateProcess(process)

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
