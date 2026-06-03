// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: No command provided.")
		fmt.Fprintln(os.Stderr, "Usage: timestamp <command> [arguments...]")
		os.Exit(1)
	}

	targetCmd := os.Args[1]
	targetArgs := os.Args[2:]

	cmd := exec.Command(targetCmd, targetArgs...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error allocating stdout pipe: %v\n", err)
		os.Exit(1)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error allocating stderr pipe: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing '%v': %v\n", cmd, err)
		os.Exit(127) // Standard UNIX exit code for command not found
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		processStream(stdoutPipe, os.Stdout)
	}()

	go func() {
		defer wg.Done()
		processStream(stderrPipe, os.Stderr)
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error harvesting process termination status: %v\n", err)
		os.Exit(1)
	}
}

func processStream(pipe io.Reader, outputStream *os.File) {
	reader := bufio.NewReader(pipe)
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			ts := time.Now().Format("15:04:05.000")
			fmt.Fprintf(outputStream, "[%s] %s", ts, line)
		}
		if err != nil {
			break // Reached EOF or pipe closed
		}
	}
}
