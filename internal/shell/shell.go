package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Result represents the result of a shell command execution
type Result struct {
	Command  string
	ExitCode int
	Stdout   string
	Stderr   string
}

// Execute runs a shell command and returns the result
func Execute(command string, args ...string) (*Result, error) {
	cmd := exec.Command(command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := &Result{
		Command: fmt.Sprintf("%s %s", command, strings.Join(args, " ")),
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		}
		return result, fmt.Errorf("command failed: %w", err)
	}

	result.ExitCode = 0
	return result, nil
}

// ExecuteInDir runs a shell command in the specified directory and returns the result
func ExecuteInDir(dir, command string, args ...string) (*Result, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := &Result{
		Command: fmt.Sprintf("%s %s", command, strings.Join(args, " ")),
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		}
		return result, fmt.Errorf("command failed: %w", err)
	}

	result.ExitCode = 0
	return result, nil
}

// CommandExists checks if a command exists in the system
func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// PrintResult prints the result of a command execution
func PrintResult(result *Result, verbose bool) {
	if verbose {
		fmt.Printf("Command: %s\n", result.Command)
		fmt.Printf("Exit Code: %d\n", result.ExitCode)
	}

	if result.Stdout != "" {
		fmt.Print(result.Stdout)
	}

	if result.Stderr != "" && verbose {
		fmt.Fprintf(os.Stderr, "Error: %s\n", result.Stderr)
	}
}
