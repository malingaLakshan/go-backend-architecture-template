//go:build windows

package cli

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/windows"
)

const (
	netstatTimeout = 5 * time.Second
)

// findListeningPID returns the PID of the process listening on the given TCP
// port. A PID of zero means no listener was found.
//
// The operating-system command and all command arguments are fixed constants.
// The port is used only to filter the returned netstat output.
func findListeningPID(port int) (int, error) {
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("invalid TCP port: %d", port)
	}

	ctx, cancel := context.WithTimeout(context.Background(), netstatTimeout)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"netstat",
		"-ano",
		"-p",
		"tcp",
	)

	output, err := cmd.Output()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return 0, fmt.Errorf("timed out while inspecting TCP listeners")
		}

		return 0, fmt.Errorf("failed to inspect TCP listeners: %w", err)
	}

	expectedPort := strconv.Itoa(port)

	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Fields(line)

		// Expected Windows netstat format:
		// TCP  local-address  remote-address  LISTENING  PID
		if len(fields) < 5 {
			continue
		}

		if !strings.EqualFold(fields[0], "TCP") {
			continue
		}

		if !strings.EqualFold(fields[3], "LISTENING") {
			continue
		}

		if portFromAddress(fields[1]) != expectedPort {
			continue
		}

		pid, err := strconv.Atoi(fields[4])
		if err != nil {
			return 0, fmt.Errorf(
				"invalid PID %q returned for port %d: %w",
				fields[4],
				port,
				err,
			)
		}

		if pid <= 0 {
			return 0, fmt.Errorf(
				"invalid PID %d returned for port %d",
				pid,
				port,
			)
		}

		return pid, nil
	}

	return 0, nil
}

// portFromAddress extracts the port from IPv4 and IPv6 netstat addresses.
//
// Examples:
//   127.0.0.1:8080       -> 8080
//   0.0.0.0:8080         -> 8080
//   [::]:8080            -> 8080
//   [::1]:8080           -> 8080
func portFromAddress(address string) string {
	index := strings.LastIndex(address, ":")
	if index < 0 || index == len(address)-1 {
		return ""
	}

	return strings.TrimSpace(address[index+1:])
}

// processEntrySize returns the actual Go structure size without importing or
// using the unsafe package.
func processEntrySize() uint32 {
	return uint32(reflect.TypeOf(windows.ProcessEntry32{}).Size())
}

// processImageName returns the executable image name for the supplied PID.
func processImageName(pid int) (string, error) {
	if pid <= 0 {
		return "", fmt.Errorf("invalid process ID: %d", pid)
	}

	snapshot, err := windows.CreateToolhelp32Snapshot(
		windows.TH32CS_SNAPPROCESS,
		0,
	)
	if err != nil {
		return "", fmt.Errorf("create process snapshot: %w", err)
	}
	defer windows.CloseHandle(snapshot)

	entry := windows.ProcessEntry32{
		Size: processEntrySize(),
	}

	if err := windows.Process32First(snapshot, &entry); err != nil {
		return "", fmt.Errorf("read first process: %w", err)
	}

	for {
		if entry.ProcessID == uint32(pid) {
			imageName := windows.UTF16ToString(entry.ExeFile[:])
			imageName = filepath.Base(strings.TrimSpace(imageName))

			if imageName == "" {
				return "", fmt.Errorf(
					"process %d has an empty executable name",
					pid,
				)
			}

			return imageName, nil
		}

		err = windows.Process32Next(snapshot, &entry)
		if err != nil {
			if errors.Is(err, windows.ERROR_NO_MORE_FILES) {
				break
			}

			return "", fmt.Errorf("read next process: %w", err)
		}
	}

	return "", fmt.Errorf("process %d was not found", pid)
}

// isReplayEngineProcess confirms that the port owner is the Replay Engine.
func isReplayEngineProcess(imageName string) bool {
	return strings.EqualFold(
		filepath.Base(strings.TrimSpace(imageName)),
		"rre.exe",
	)
}

// terminateProcessTree terminates the child processes first and then the root
// Replay Engine process. It uses the Windows API and does not execute taskkill
// or another operating-system command with dynamic input.
func terminateProcessTree(pid int) error {
	if pid <= 0 {
		return fmt.Errorf("invalid process ID: %d", pid)
	}

	processIDs, err := processTreePIDs(uint32(pid))
	if err != nil {
		return fmt.Errorf(
			"resolve process tree for PID %d: %w",
			pid,
			err,
		)
	}

	var terminationErrors []error

	// processTreePIDs returns child processes before their parent.
	for _, processID := range processIDs {
		if processID == 0 {
			continue
		}

		if err := terminateWindowsProcess(processID); err != nil {
			terminationErrors = append(
				terminationErrors,
				fmt.Errorf(
					"terminate PID %d: %w",
					processID,
					err,
				),
			)
		}
	}

	if len(terminationErrors) > 0 {
		return errors.Join(terminationErrors...)
	}

	return nil
}

// processTreePIDs returns all descendants followed by the root PID. This order
// allows child processes to be terminated before their parent.
func processTreePIDs(rootPID uint32) ([]uint32, error) {
	if rootPID == 0 {
		return nil, fmt.Errorf("root process ID must be greater than zero")
	}

	snapshot, err := windows.CreateToolhelp32Snapshot(
		windows.TH32CS_SNAPPROCESS,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("create process snapshot: %w", err)
	}
	defer windows.CloseHandle(snapshot)

	childrenByParent := make(map[uint32][]uint32)

	entry := windows.ProcessEntry32{
		Size: processEntrySize(),
	}

	if err := windows.Process32First(snapshot, &entry); err != nil {
		return nil, fmt.Errorf("read first process: %w", err)
	}

	for {
		if entry.ProcessID != 0 &&
			entry.ProcessID != entry.ParentProcessID {

			childrenByParent[entry.ParentProcessID] = append(
				childrenByParent[entry.ParentProcessID],
				entry.ProcessID,
			)
		}

		err = windows.Process32Next(snapshot, &entry)
		if err != nil {
			if errors.Is(err, windows.ERROR_NO_MORE_FILES) {
				break
			}

			return nil, fmt.Errorf("read next process: %w", err)
		}
	}

	result := make([]uint32, 0)
	visited := make(map[uint32]bool)

	var collectChildren func(uint32)

	collectChildren = func(parentPID uint32) {
		if visited[parentPID] {
			return
		}

		visited[parentPID] = true

		for _, childPID := range childrenByParent[parentPID] {
			if childPID == 0 || childPID == parentPID {
				continue
			}

			collectChildren(childPID)
		}

		result = append(result, parentPID)
	}

	collectChildren(rootPID)

	return result, nil
}

// terminateWindowsProcess terminates one process using the Windows API.
func terminateWindowsProcess(pid uint32) error {
	if pid == 0 {
		return fmt.Errorf("process ID must be greater than zero")
	}

	processHandle, err := windows.OpenProcess(
		windows.PROCESS_TERMINATE,
		false,
		pid,
	)
	if err != nil {
		// The process may have exited after the snapshot was created.
		if errors.Is(err, windows.ERROR_INVALID_PARAMETER) {
			return nil
		}

		return fmt.Errorf("open process: %w", err)
	}
	defer windows.CloseHandle(processHandle)

	if err := windows.TerminateProcess(processHandle, 1); err != nil {
		// Treat an already-terminated process as successful.
		if errors.Is(err, windows.ERROR_INVALID_HANDLE) ||
			errors.Is(err, windows.ERROR_INVALID_PARAMETER) {

			return nil
		}

		return fmt.Errorf("terminate process: %w", err)
	}

	return nil
}