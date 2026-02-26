---
name: authcurl
description: This skill should be used when the user asks to "use authcurl", "configure authcurl", "set up curl authentication", "inject auth headers into curl", "add API token to curl", or mentions authcurl, authenticated curl requests, or curl with automatic auth injection.
version: 0.1.0
---

# authcurl

A drop-in `curl` wrapper that auto-injects authentication headers based on the request URL. All curl args, exit codes, signals, and stdin/stdout pass through directly via `syscall.Exec`.

## Installation

```bash
# From source
go install github.com/nguyenvanduocit/authcurl@latest

# Homebrew
brew install nguyenvanduocit/tap/authcurl

# Or build from repo
git clone https://github.com/nguyenvanduocit/authcurl
cd authcurl && go build -o authcurl .
```

## Setup

```bash
# Create config with examples
authcurl init

# Edit the config
$EDITOR ~/.config/authcurl/config.yaml

# Verify profiles
authcurl list
```

## Config Format

Config location: `~/.config/authcurl/config.yaml`

Profiles are matched top-to-bottom against the full request URL. The first match wins.

```yaml
profiles:
  - name: github
    match: "https://api.github.com/**"
    bearer: "${GITHUB_TOKEN}"
```

Each profile requires:
- `name` — identifier shown in stderr on match
- `match` — URL pattern (`**` or `*` suffix for prefix match, otherwise exact)
- One auth field: `bearer`, `headers`, or `basic_auth`

## Auth Types

### Bearer Token

Injects `-H "Authorization: Bearer <token>"`.

```yaml
- name: github
  match: "https://api.github.com/**"
  bearer: "${GITHUB_TOKEN}"
```

### Custom Headers

Injects `-H "Key: Value"` for each entry.

```yaml
- name: internal-api
  match: "https://api.internal.example.com/**"
  headers:
    X-API-Key: "${INTERNAL_API_KEY}"
    X-Service: "my-service"
```

### Basic Auth

Injects `-u username:password`.

```yaml
- name: legacy
  match: "https://legacy.example.com/**"
  basic_auth:
    username: "${LEGACY_USER}"
    password: "${LEGACY_PASS}"
```

## Environment Variables

All config values support `${VAR_NAME}` expansion at runtime. Always use env vars for secrets — never hardcode tokens.

Auth type detection happens before env var expansion, so `bearer: "${TOKEN}"` is treated as bearer even if `TOKEN` is unset.

## Commands

| Command | Description |
|---|---|
| `authcurl [curl-args...] <url>` | Run curl with auto-injected auth |
| `authcurl init` | Create `~/.config/authcurl/config.yaml` with examples |
| `authcurl list` | Print configured profiles with match patterns and auth types |

## Additional Resources

### Reference Files

For detailed patterns and troubleshooting, consult:
- **`references/patterns.md`** — Common config patterns (GitHub, AWS, multi-environment, private registries)
- **`references/troubleshooting.md`** — Diagnosing config, matching, and auth issues
