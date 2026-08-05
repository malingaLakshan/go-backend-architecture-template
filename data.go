// Run replays all RawReads with real-time pacing.
// It respects context cancellation for Ctrl+C handling.
func (s *Service) Run(
	ctx context.Context,
	reads []recording.RawRead,
) *Status {
	status := NewStatus(len(reads))
	status.StartTime = time.Now()

	s.Logger.Info("Playback started, total records: %d", len(reads))

	// Open replay output file.
	outputPath := filepath.Join(replayOutputDir, replayOutputFile)

	if err := os.MkdirAll(replayOutputDir, 0o755); err != nil {
		s.Logger.Error("Failed to create output directory: %v", err)
		status.State = "failed"
		status.EndTime = time.Now()
		return status
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		s.Logger.Error("Failed to create replay output file: %v", err)
		status.State = "failed"
		status.EndTime = time.Now()
		return status
	}
	defer outputFile.Close()

	// Anchor all records to one absolute replay timeline.
	// This prevents payload processing and HTTP request time from
	// accumulating between consecutive records.
	replayStart := time.Now()

	var firstInjectionTime time.Time
	if len(reads) > 0 {
		firstInjectionTime = reads[0].InjectionTime
	}

	for i, read := range reads {
		// Check for cancellation before processing each record.
		select {
		case <-ctx.Done():
			s.Logger.Warn("Playback aborted by user")
			status.State = "aborted"
			status.EndTime = time.Now()
			return status
		default:
		}

		isFirst := i == 0

		// Calculate the absolute scheduled time of this record.
		targetOffset := read.InjectionTime.Sub(firstInjectionTime)
		targetTime := replayStart.Add(targetOffset)
		delay := time.Until(targetTime)

		// Wait only for the remaining time.
		// If playback is already late, send immediately to catch up.
		if delay > 0 {
			timer := time.NewTimer(delay)

			select {
			case <-ctx.Done():
				timer.Stop()
				s.Logger.Warn(
					"Playback aborted by user during pacing wait",
				)
				status.State = "aborted"
				status.EndTime = time.Now()
				return status

			case <-timer.C:
			}
		}

		// Build the ProtoReaderBundle payload from the RawRead.
		payload, err := BuildPayload(&read, s.SiteID)
		if err != nil {
			s.Logger.Error(
				"Failed to build payload for ReadID %s: %v",
				read.ReadID,
				err,
			)
			status.RecordFailure()
			continue
		}

		// Serialize the payload for the output file and sending.
		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			s.Logger.Error(
				"Failed to marshal payload for ReadID %s: %v",
				read.ReadID,
				err,
			)
			status.RecordFailure()
			continue
		}

		// Write to replay output file.
		fmt.Fprintln(outputFile, string(payloadJSON))

		// Print the first payload to the terminal.
		if isFirst {
			prettyJSON, _ := json.MarshalIndent(payload, "", "  ")

			fmt.Println()
			fmt.Println("First replay output payload:")
			fmt.Println(string(prettyJSON))
			fmt.Println()
		}

		// Send to target.
		if err := s.Injector.Send(payload); err != nil {
			s.Logger.Error(
				"Failed to send ReadID %s at %s: %v",
				read.ReadID,
				read.InjectionTime.Format(time.RFC3339Nano),
				err,
			)
			status.RecordFailure()
			continue
		}

		status.RecordSuccess()
	}

	status.State = "completed"
	if status.Failed > 0 {
		status.State = "completed_with_errors"
	}

	status.EndTime = time.Now()

	s.Logger.Info(
		"Playback finished, successful: %d, failed: %d",
		status.Successful,
		status.Failed,
	)

	return status
}