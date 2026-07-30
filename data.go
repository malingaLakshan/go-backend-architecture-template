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
		Size: uint32(unsafe.Sizeof(windows.ProcessEntry32{})),
	}

	err = windows.Process32First(snapshot, &entry)
	if err != nil {
		return "", fmt.Errorf("read first process: %w", err)
	}

	for {
		if entry.ProcessID == uint32(pid) {
			imageName := windows.UTF16ToString(entry.ExeFile[:])
			return filepath.Base(imageName), nil
	}

		err = windows.Process32Next(snapshot, &entry)
		if err != nil {
			if errors.Is(err, syscall.ERROR_NO_MORE_FILES) {
				break
			}

			return "", fmt.Errorf("read next process: %w", err)
		}
	}

	return "", fmt.Errorf("process %d was not found", pid)
}