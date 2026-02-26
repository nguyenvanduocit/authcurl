# Common Config Patterns

## GitHub API

```yaml
- name: github
  match: "https://api.github.com/**"
  bearer: "${GITHUB_TOKEN}"
```
```bash
authcurl https://api.github.com/user
authcurl https://api.github.com/repos/owner/repo/issues
authcurl -X POST -d '{"title":"bug"}' https://api.github.com/repos/owner/repo/issues
```

## AWS API Gateway

```yaml
- name: aws-gateway
  match: "https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/**"
  headers:
    x-api-key: "${AWS_API_KEY}"
```

## Multiple Environments

```yaml
- name: prod-api
  match: "https://api.example.com/**"
  bearer: "${PROD_TOKEN}"

- name: staging-api
  match: "https://staging-api.example.com/**"
  bearer: "${STAGING_TOKEN}"
```

Profile order matters — place more specific patterns before broader ones.

## Private Package Registries

```yaml
- name: npm-private
  match: "https://registry.npmjs.org/**"
  headers:
    Authorization: "Bearer ${NPM_TOKEN}"
```

## Datadog API

```yaml
- name: datadog
  match: "https://api.datadoghq.com/**"
  headers:
    DD-API-KEY: "${DD_API_KEY}"
    DD-APPLICATION-KEY: "${DD_APP_KEY}"
```

## URL Pattern Matching

| Pattern | Matches |
|---|---|
| `https://api.github.com/**` | Any URL starting with `https://api.github.com/` |
| `https://api.github.com/*` | Same — both `*` and `**` perform prefix matching |
| `https://api.github.com/user` | Exact URL only |

Both `*` and `**` suffixes strip the wildcard and check `strings.HasPrefix`. Without a wildcard, the pattern must match the URL exactly.
