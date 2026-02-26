# Troubleshooting

## "config not found at ~/.config/authcurl/config.yaml"

Run `authcurl init` to create the config file, then edit it with profiles.

## Profile not matching — no auth injected

- Run `authcurl list` to verify configured patterns
- Patterns use simple prefix matching: `**` and `*` both strip the suffix and check `strings.HasPrefix`
- The URL must start with the prefix literally, including scheme (`https://`)
- Only the first matching profile is used — check profile order

## Empty token injected (request fails with 401)

- The env var referenced in the config is unset or empty
- Auth type detection uses the literal `${VAR}` string (before expansion), so the profile still matches
- Verify the variable is exported: `echo $YOUR_VAR`

## "curl not found in PATH"

Install curl: `brew install curl` (macOS) or `apt install curl` (Linux).

## "config already exists" error on init

`authcurl init` refuses to overwrite existing config. Edit `~/.config/authcurl/config.yaml` directly.
