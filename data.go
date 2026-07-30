// terminateProcessTree terminates the supplied process and its child processes.
func terminateProcessTree(pid int) error {
	if pid <= 0 {
		return fmt.Errorf("invalid process ID: %d", pid)
	}

	if pid == os.Getpid() {
		return fmt.Errorf("refusing to terminate the current process")
	}

	processIDs, err := processTreePIDs(uint32(pid))
	if err != nil {
		return err
	}

	var terminationErrors []error

	// processTreePIDs returns child processes before the parent.
	for _, processID := range processIDs {
		if processID == uint32(os.Getpid()) {
			continue
		}

		err := terminateWindowsProcess(processID)
		if err != nil {
			terminationErrors = append(
				terminationErrors,
				fmt.Errorf("terminate process %d: %w", processID, err),
			)
		}
	}

	if len(terminationErrors) > 0 {
		return errors.Join(terminationErrors...)
	}

	return nil
}// processTreePIDs returns all descendants followed by the root process.
func processTreePIDs(rootPID uint32) ([]uint32, error) {
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
		Size: uint32(unsafe.Sizeof(windows.ProcessEntry32{})),
	}

	err = windows.Process32First(snapshot, &entry)
	if err != nil {
		return nil, fmt.Errorf("read first process: %w", err)
	}

	for {
		childrenByParent[entry.ParentProcessID] = append(
			childrenByParent[entry.ParentProcessID],
			entry.ProcessID,
		)

		err = windows.Process32Next(snapshot, &entry)
		if err != nil {
			if errors.Is(err, syscall.ERROR_NO_MORE_FILES) {
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
	processHandle, err := windows.OpenProcess(
		windows.PROCESS_TERMINATE,
		false,
		pid,
	)
	if err != nil {
		// The process may have already exited while the process tree was read.
		if errors.Is(err, syscall.ERROR_INVALID_PARAMETER) {
			return nil
		}

		return fmt.Errorf("open process: %w", err)
	}
	defer windows.CloseHandle(processHandle)

	if err := windows.TerminateProcess(processHandle, 1); err != nil {
		return fmt.Errorf("terminate process: %w", err)
	}

	return nil
}

go get golang.org/x/sys/windows
go mod tidygofmt -w internal/cli/stop_windows.go
go test ./...
go vet ./...






