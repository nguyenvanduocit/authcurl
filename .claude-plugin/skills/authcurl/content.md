# authcurl

**Triggers:** authcurl, curl auth, curl authentication, API auth headers, bearer token curl, inject auth header, authenticated curl request

`authcurl` is a drop-in curl wrapper that auto-injects authentication headers based on the request URL. Use it anywhere you'd use curl — it passes all args, exit codes, signals, stdin/stdout directly through to curl via `syscall.Exec`.

## Installation

```bash
# From source
go install github.com/nguyenvanduocit/authcurl@latest

# Or build from repo
git clone https://github.com/nguyenvanduocit/authcurl
cd authcurl && go build -o authcurl . && mv authcurl /usr/local/bin/
```

## Quick Start

```bash
# 1. Create config with examples
authcurl init

# 2. Edit ~/.config/authcurl/config.yaml to add your profiles

# 3. Use exactly like curl
authcurl https://api.github.com/user
authcurl -X POST -d '{"key":"val"}' https://api.github.com/repos/owner/repo/issues
```

## Config File

Location: `~/.config/authcurl/config.yaml`

```yaml
profiles:
  - name: github
    match: "https://api.github.com/**"
    bearer: "${GITHUB_TOKEN}"
```

Each profile has:
- `name` — label shown in stderr when matched (`authcurl: using profile "github"`)
- `match` — URL glob pattern (matched against the full URL)
- one auth field: `bearer`, `headers`, or `basic_auth`

## URL Pattern Matching

Patterns are matched against the full request URL in profile order — first match wins.

| Pattern | Matches |
|---|---|
| `https://api.github.com/**` | any URL starting with `https://api.github.com/` |
| `https://api.github.com/*` | same (both `*` and `**` use prefix matching) |
| `https://api.github.com/user` | exact URL only |

## Auth Types

### 1. Bearer Token

Injects `-H "Authorization: Bearer <token>"`.

```yaml
- name: github
  match: "https://api.github.com/**"
  bearer: "${GITHUB_TOKEN}"
```

### 2. Custom Headers

Injects `-H "Key: Value"` for each entry.

```yaml
- name: internal-api
  match: "https://api.internal.example.com/**"
  headers:
    X-API-Key: "${INTERNAL_API_KEY}"
    X-Service: "my-service"
```

### 3. Basic Auth

Injects `-u username:password`.

```yaml
- name: legacy
  match: "https://legacy.example.com/**"
  basic_auth:
    username: "${LEGACY_USER}"
    password: "${LEGACY_PASS}"
```

## Environment Variables

All config values support `${VAR_NAME}` expansion. Always use env vars for secrets — never hardcode tokens in the config file.

```bash
# Set before running authcurl, or export in your shell profile
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
export INTERNAL_API_KEY="sk-xxxxxxxx"
```

Auth type detection happens before env var expansion, so a profile with `bearer: "${TOKEN}"` is always treated as bearer even if `TOKEN` is unset at detection time.

## Common Patterns

**GitHub API**
```yaml
- name: github
  match: "https://api.github.com/**"
  bearer: "${GITHUB_TOKEN}"
```
```bash
authcurl https://api.github.com/user
authcurl https://api.github.com/repos/owner/repo/issues
```

**AWS API Gateway / custom internal APIs**
```yaml
- name: aws-gateway
  match: "https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/**"
  headers:
    x-api-key: "${AWS_API_KEY}"
```

**Multiple environments with different keys**
```yaml
- name: prod-api
  match: "https://api.example.com/**"
  bearer: "${PROD_TOKEN}"

- name: staging-api
  match: "https://staging-api.example.com/**"
  bearer: "${STAGING_TOKEN}"
```

**NPM registry / package feeds**
```yaml
- name: npm-private
  match: "https://registry.npmjs.org/**"
  headers:
    Authorization: "Bearer ${NPM_TOKEN}"
```

## Commands

| Command | Description |
|---|---|
| `authcurl [curl-args...] <url>` | Run curl with auto-injected auth |
| `authcurl init` | Create `~/.config/authcurl/config.yaml` with examples |
| `authcurl list` | Print all configured profiles with their match patterns and auth types |

## Troubleshooting

**"config not found at ~/.config/authcurl/config.yaml"**
Run `authcurl init` to create the config, then edit it.

**Profile not matching — no auth injected**
- Run `authcurl list` to see configured patterns
- Patterns use simple prefix matching: `**` and `*` both strip the suffix and check `strings.HasPrefix`
- The URL must start with the prefix literally, including scheme (`https://`)
- Only the first matching profile is used

**Empty token injected (request fails with 401)**
- The env var is unset — auth type was detected from the literal `${VAR}` string, but the value expanded to empty
- Run `echo $YOUR_VAR` to confirm the variable is exported in the current shell

**"curl not found in PATH"**
Install curl: `brew install curl` (macOS) or `apt install curl` (Linux).

**Config already exists error on init**
`authcurl init` refuses to overwrite an existing config. Edit `~/.config/authcurl/config.yaml` directly.
