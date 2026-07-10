README was not updated. Please update `resonate-replay-engine/README.md` for the new versioning feature.

Add a clear section for CLI versioning.

Required README updates:
1. Add version command usage:
   - `.\rre.exe version`
   - `.\rre.exe --version`
   - `.\rre.exe -v`

2. Add expected output:
   `rre version dev`
   for normal local build.

3. Add build command with injected version:
   `go build -ldflags="-X main.version=v0.1.0" -o rre.exe ./cmd/rre`

4. Add expected output after version-injected build:
   `rre version v0.1.0`

5. Add a short note:
   - `dev` is the default version for local builds.
   - Release builds should inject the release tag using `-ldflags`.
   - Git tags should follow the monorepo tool format: `resonate-replay-engine/vMAJOR.MINOR.PATCH`

6. Do not change unrelated docs.
7. Do not add generated files like `rre.exe`, logs, sqlite wal/shm, or bin files.

After update, show the README changes only.