Review the replay pacing logic in internal/replay/service.go.

The replay currently waits using the timestamp difference between consecutive InjectionTime values. Payload creation, JSON writing, and HTTP sending happen after each wait, causing cumulative timing drift.

Example:

* Recorded: 7045 RawReads
* Original duration: 11m42.9801641s
* Replay duration: 12m14.4s
* Drift: ~31.4s

Implement the smallest safe change so replay follows the original recording timeline using an absolute schedule from the first InjectionTime, preventing accumulated drift.

Keep existing logging, cancellation, status updates, error handling, public APIs, and output behavior unchanged. Do not refactor unrelated code.

First explain the approach briefly, then provide only the code changes.