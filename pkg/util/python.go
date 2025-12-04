package util

import (
	"bytes"
	"fmt"
	"os/exec"
)

const (
	pythonExecutable = `C:\ProgramData\Chocolatey\bin\python3.14.exe`
	// pythonExecutable = "python3"
)

// RunPythonScript executes a Python script and returns its stdout as a string.
// If the script returns a non-zero exit code, stderr is included in the error.
func RunPythonScript(scriptPath string, args ...string) (string, error) {
	// Build full command: python3 script.py arg1 arg2 ...
	cmdArgs := append([]string{scriptPath}, args...)
	cmd := exec.Command(pythonExecutable, cmdArgs...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Attach stderr to the error message for easier debugging
		return "", fmt.Errorf(
			"python script failed: %w\nstderr: %s",
			err, stderr.String(),
		)
	}

	return stdout.String(), nil
}
