# AGENTS.md

`protect` is a Go CLI and Bubble Tea TUI for UniFi Protect: switch viewports
between liveviews and move PTZ cameras to presets.

## Layout

- `main.go` - entry point; calls `cmd.Execute()`.
- `cmd/root.go` - single Cobra root command, flag-driven (no subcommands).
  `PersistentPreRunE` loads config, applies `--url`/`--token`/`--log-level`
  overrides, sets the log level and validates (skipped for `--version`).
- `internal/config` - Viper singleton. Searches `$XDG_CONFIG_HOME/protect`,
  `~/.config/protect`, then `os.UserConfigDir()/protect`. Env prefix `PROTECT_`.
- `internal/client` - HTTP client for `/proxy/protect/integration/v1/*`,
  authenticated with the `X-API-Key` header.
- `internal/logger` - global zap logger, disabled (`none`) by default.
- `internal/tui` - Bubble Tea state machine: main menu -> viewports ->
  liveviews, and main menu -> cameras -> presets (-1 to 9).

## Naming gotchas

- `Viewport` is an alias of `Viewer` (API `/viewers`).
- `Camera` is an alias of `Liveview`; `ListCameras()` returns **liveviews**.
- `ListPTZCameras()` filters `/cameras` with `PTZCamera.HasPTZ()`. The
  integration API has no PTZ capability flag; this checks the undocumented
  `type` field (model name, e.g. "UVC G6 PTZ") for "PTZ".

## Development

Tasks follow the `methridge/taskfiles` standard. Never hand-edit
`.taskfiles/shared/` (`task sync` overwrites it); repo tasks live in
`.taskfiles/project/project.yml`.

```bash
task build   # bin/protect
task test    # go test -v ./...
task lint    # go fmt + go vet
```

## Releases

GoReleaser runs on `v*` tags (`.github/workflows/release.yml`) and publishes a
Homebrew cask to `methridge/homebrew-tap`. Version info is injected with
`-X github.com/methridge/protect/cmd.{version,commit,date}`.

## Conventions

- Work on a `<type>/<slug>` branch and merge via PR; `main` is protected by the
  `no-commit-to-branch` pre-commit hook.
- Conventional Commit messages.
