//go:build !windows

package cli

import "fmt"

func findListeningPID(port int) (int, error) {
	return 0, fmt.Errorf(
		"the stop command is currently supported only on Windows",
	)
}

func processImageName(pid int) (string, error) {
	return "", fmt.Errorf(
		"the stop command is currently supported only on Windows",
	)
}

func isReplayEngineProcess(imageName string) bool {
	return false
}

func terminateProcessTree(pid int) error {
	return fmt.Errorf(
		"the stop command is currently supported only on Windows",
	)
}