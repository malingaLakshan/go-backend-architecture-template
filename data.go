The Play request reaches POST /api/v1/playback/start with HTTP 200, but the response Content-Type is text/html instead of JSON.

Investigate only the UI API routing.

Find why /api/v1/playback/start is being handled by the static file server (index.html) instead of the playback handler.

Ensure:
- all /api/* routes are registered before the static handler
- POST /api/v1/playback/start executes the existing replay handler
- successful responses return application/json
- static files only serve non-API routes

Do not change replay logic or validation.
Show the root cause before modifying code.