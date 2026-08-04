Fix only the currently failing Go unit tests shown by go test ./... -count=1.

Failing areas:

* internal/mocktarget
    * approved SiteGraph directory path tests
    * site store summary counts for readers and antenna ports
* internal/recording
    * RawReads test fixtures fail scanning reader_id because test data type does not match the repository model
* internal/replay
    * HTTPS URL tests fail because the injector currently accepts only HTTP

Requirements:

* First identify whether each failure is caused by outdated tests or an actual implementation bug.
* Preserve current production behavior unless the test clearly exposes a real bug.
* Update tests when expectations or fixtures are outdated.
* Do not weaken path-security validation.
* Do not add broad refactors or unrelated changes.
* Keep changes limited to the failing packages.
* Run the affected package tests first, then go test ./... -count=1.
* Summarize exactly which files changed and why.

Do not modify generated files, logging behavior, CLI behavior, or public APIs unless required to fix a confirmed defect.