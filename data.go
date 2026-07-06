Implemented Replay Engine CLI MVP features including sample recording, summary, mock server, validation, replay flow, pacing, output logging, and Cycode ## Summary

Fixed Cycode security scan issues in the Replay Engine CLI MVP.

## Changes

- Fixed SQL injection warning by using allowed hardcoded SQL queries.
- Fixed SSRF warning by validating the target URL.
- Restricted replay target to the local mock server for MVP.
- Updated injector to use safe allowlisted endpoints.

## Testing

- Ran gofmt.
- Ran go test ./...
- Cycode scan passed.



