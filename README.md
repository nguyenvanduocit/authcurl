# authcurl

A `curl` wrapper that automatically injects authentication credentials based on URL pattern matching. No more copy-pasting tokens into every request.

```
authcurl https://api.github.com/user
# authcurl: using profile "github"
# ... curl executes with Authorization: Bearer <token> injected
```


## Installation

**From source:**

```sh
go install github.com/nguyenvanduocit/authcurl@latest
```

**Binary releases:**

Download pre-built binaries from the [releases page](https://github.com/nguyenvanduocit/authcurl/releases).

**Homebrew:**

```sh
# Coming soon
brew install nguyenvanduocit/tap/authcurl
```


## Quick Start

```sh
# 1. Create the config file
authcurl init

# 2. Edit it with your credentials
$EDITOR ~/.config/authcurl/config.yaml

# 3. Make requests — auth is injected automatically
authcurl https://api.github.com/user
authcurl -X POST -d '{"title":"bug"}' https://api.github.com/repos/owner/repo/issues
```


## Config File

Location: `~/.config/authcurl/config.yaml`

Profiles are matched top-to-bottom. The first matching profile wins.

```yaml
profiles:
  - name: github
    match: "https://api.github.com/**"
    bearer: "${GITHUB_TOKEN}"

  - name: internal-api
    match: "https://api.example.com/**"
    headers:
      X-API-Key: "${API_KEY}"
      X-Custom: "some-value"

  - name: legacy
    match: "https://legacy.example.com/**"
    basic_auth:
      username: "${LEGACY_USER}"
      password: "${LEGACY_PASS}"
```

### Auth Types

| Type | Config key | Injected as |
|------|-----------|-------------|
| Bearer token | `bearer` | `-H "Authorization: Bearer <token>"` |
| Custom headers | `headers` | `-H "Key: Value"` for each entry |
| Basic auth | `basic_auth` | `-u username:password` |

### URL Pattern Matching

| Pattern | Matches |
|---------|---------|
| `https://api.github.com/**` | Any URL starting with that prefix |
| `https://api.example.com/*` | Any URL starting with that prefix |
| `https://api.example.com/v1/users` | Exact URL only |

Both `**` and `*` suffixes perform prefix matching. Without a wildcard suffix, the pattern must match the URL exactly.

### Environment Variables

All credential values support `${VAR_NAME}` expansion. Values are resolved from the environment at runtime. This keeps secrets out of the config file.

```yaml
bearer: "${GITHUB_TOKEN}"
headers:
  X-API-Key: "${MY_API_KEY}"
basic_auth:
  username: "${DB_USER}"
  password: "${DB_PASS}"
```


## Examples

**GitHub API:**

```sh
export GITHUB_TOKEN=ghp_...

authcurl https://api.github.com/user
authcurl https://api.github.com/repos/owner/repo/issues
authcurl -X POST -H "Content-Type: application/json" \
  -d '{"title":"New issue"}' \
  https://api.github.com/repos/owner/repo/issues
```

**AWS (custom auth header):**

```yaml
- name: aws-internal
  match: "https://my-service.us-east-1.amazonaws.com/**"
  headers:
    Authorization: "${AWS_AUTH_HEADER}"
    X-Amz-Security-Token: "${AWS_SESSION_TOKEN}"
```

**Private API with API key:**

```yaml
- name: datadog
  match: "https://api.datadoghq.com/**"
  headers:
    DD-API-KEY: "${DD_API_KEY}"
    DD-APPLICATION-KEY: "${DD_APP_KEY}"
```

**Multiple environments:**

```yaml
profiles:
  - name: prod
    match: "https://api.myapp.com/**"
    bearer: "${PROD_TOKEN}"

  - name: staging
    match: "https://staging-api.myapp.com/**"
    bearer: "${STAGING_TOKEN}"
```


## Commands

```
authcurl [curl-options...] <url>    Execute curl with auto-injected auth
authcurl init                       Create config file at ~/.config/authcurl/config.yaml
authcurl list                       List configured profiles and their auth types
```


## How It Works

1. authcurl parses your arguments to extract the request URL (bare `https://` argument or `--url=` flag).
2. It walks the profiles list and finds the first profile whose `match` pattern fits the URL.
3. Auth flags (`-H`, `-u`) are prepended to the argument list.
4. `syscall.Exec` replaces the authcurl process with `curl`. There is no wrapper process — exit codes, signals, stdin, stdout, and stderr all belong to curl directly.

If no profile matches, curl runs with your original arguments unchanged.


## Requirements

- Go 1.21+ (to build from source)
- `curl` available in `PATH`


## License

MIT
