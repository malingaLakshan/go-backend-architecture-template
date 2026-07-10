Please add config-file support to the summary command without breaking the existing -file behavior.

Current behavior:
rre summary -file data/recording_001.sqlite

New behavior needed:
rre summary -config configs/pass_config.json

Rules:
1. If -config is provided, load RunConfig and use recording_file as the summary input file.
2. If -file is provided, keep existing behavior.
3. If both -config and -file are provided, let -config take priority or return a clear error. Prefer clear error.
4. Summary does not need target_url, site_id, mock_mode, or mock_site_file.
5. Do not change validate/play/mock-server behavior.
6. Add/update help text to show both summary examples.
7. Run gofmt, go test ./..., and go build -o rre.exe ./cmd/rre.