# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

authcurl is a lightweight `curl` wrapper in Go that automatically injects authentication credentials based on URL pattern matching. It uses `syscall.Exec` to replace the process with curl directly — zero wrapper overhead.

## Build & Run

```sh
go build -o authcurl .       # Build binary
go run . <url>               # Run directly
go test ./...                # Run tests (none yet)
```

## CI

GitHub Actions runs on push to `main` and on PRs:
- `go build ./...`
- `golangci-lint` (latest)
- `go test ./...`

Releases use GoReleaser (`.goreleaser.yaml`) triggered by tags, builds for linux/darwin on amd64/arm64.

## Architecture

Three files, single `main` package:

- **main.go** — Entry point, command routing (`init`, `list`, or curl exec)
- **config.go** — YAML config loading (`~/.config/authcurl/config.yaml`), profile matching, env var expansion
- **curl.go** — URL extraction from args, auth header injection, `X-Powered-By` header, `syscall.Exec` to curl

Flow: parse args → extract URL → match profile → inject auth flags → exec curl.

## Key Patterns

- **Auth types**: bearer token, custom headers, basic auth
- **URL matching**: prefix match with `**`/`*` suffix, or exact match (no wildcard)
- **Env vars**: `${VAR_NAME}` syntax expanded at runtime via `os.ExpandEnv`
- **Process replacement**: `syscall.Exec` replaces authcurl with curl (exit codes, signals, stdio pass through)

## Conventions

- Keep it minimal — single binary, no subpackages, no external deps beyond `yaml.v3`
- Config lives at `~/.config/authcurl/config.yaml`
- Errors go to stderr with `authcurl:` prefix
- Profile match is first-match-wins, top to bottom
