//go:build windows

package cli

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// findListeningPID returns the PID of the process listening on the given TCP port.
func findListeningPID(port int) (int, error) {
	cmd := exec.Command(
		"netstat",
		"-ano",
		"-p",
		"tcp",
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to inspect listening ports: %w", err)
	}

	expectedSuffix := ":" + strconv.Itoa(port)

	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Fields(line)

		// Expected Windows netstat format:
		// TCP local-address remote-address LISTENING pid
		if len(fields) < 5 {
			continue
		}

		protocol := strings.ToUpper(fields[0])
		localAddress := fields[1]
		state := strings.ToUpper(fields[3])
		pidText := fields[4]

		if protocol != "TCP" {
			continue
		}

		if state != "LISTENING" {
			continue
		}

		if !strings.HasSuffix(localAddress, expectedSuffix) {
			continue
		}

		pid, err := strconv.Atoi(pidText)
		if err != nil {
			return 0, fmt.Errorf(
				"invalid PID %q returned for port %d: %w",
				pidText,
				port,
				err,
			)
		}

		return pid, nil
	}

	return 0, nil
}

// processImageName returns the executable image name for the supplied PID.
func processImageName(pid int) (string, error) {
	cmd := exec.Command(
		"tasklist",
		"/FI",
		fmt.Sprintf("PID eq %d", pid),
		"/FO",
		"CSV",
		"/NH",
	)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf(
			"failed to inspect process %d: %w",
			pid,
			err,
		)
	}

	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" ||
		strings.HasPrefix(strings.ToUpper(trimmed), "INFO:") {

		return "", fmt.Errorf("process %d was not found", pid)
	}

	reader := csv.NewReader(bytes.NewBufferString(trimmed))

	record, err := reader.Read()
	if err != nil {
		return "", fmt.Errorf(
			"failed to parse process information for PID %d: %w",
			pid,
			err,
		)
	}

	if len(record) == 0 {
		return "", fmt.Errorf(
			"process information for PID %d was empty",
			pid,
		)
	}

	return strings.TrimSpace(record[0]), nil
}

// isReplayEngineProcess confirms that the port owner is the Replay Engine.
func isReplayEngineProcess(imageName string) bool {
	return strings.EqualFold(
		strings.TrimSpace(imageName),
		"rre.exe",
	)
}

// terminateProcessTree terminates the process and any child processes.
func terminateProcessTree(pid int) error {
	cmd := exec.Command(
		"taskkill",
		"/PID",
		strconv.Itoa(pid),
		"/T",
		"/F",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"failed to terminate process %d: %s: %w",
			pid,
			strings.TrimSpace(string(output)),
			err,
		)
	}

	return nil
}