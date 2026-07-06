Cloud OAuth credentials (gitignored)

Place OAuth JSON files here for the unified Sylos server. The API reads
these at runtime and performs OAuth token exchange server-side — client
secrets never reach the browser.

  creds/google-oauth.json
  creds/dropbox-oauth.json

Use the same JSON format as Sylos-UI (Google "installed"/"web" wrapper
or flat client_id/client_secret objects).

If files are not here, the server also checks ui/creds/ (the UI submodule)
as a fallback.

Default location is ./creds next to config.yaml. Override in config.yaml:

  runtime:
    oauth_creds_dir: "./creds"

Restart the Sylos server after changing credential files.

Do not commit real credentials.
