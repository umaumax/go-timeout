package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"time"
)

func execCommand(command string, args []string, duration time.Duration) (exitCode int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return
	}
	err = cmd.Wait()
	exitCode = cmd.ProcessState.ExitCode()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			if ctx_err := ctx.Err(); ctx_err != nil {
				log.Printf("timeout: %v\n", err)
			}
			err = nil
		}
	}
	return
}

var (
	t time.Duration
)

func init() {
	flag.DurationVar(&t, "t", time.Duration(0), "timeout")
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return
	}
	command := args[0]
	args = args[1:]

	exitCode, err := execCommand(command, args, t)
	if err != nil {
		log.Fatalf("exec command error:%s\n", err)
	}
	os.Exit(exitCode)
}
