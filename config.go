package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Profiles []Profile `yaml:"profiles"`
}

type Profile struct {
	Name      string            `yaml:"name"`
	Match     string            `yaml:"match"`
	Bearer    string            `yaml:"bearer,omitempty"`
	Headers   map[string]string `yaml:"headers,omitempty"`
	BasicAuth *BasicAuth        `yaml:"basic_auth,omitempty"`
	authType  string
}

type BasicAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "authcurl", "config.yaml")
}

func loadConfig() (*Config, error) {
	path := configPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config not found at %s", path)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		// Detect auth type before expansion (env vars may resolve to empty)
		switch {
		case p.Bearer != "":
			p.authType = "bearer"
		case p.BasicAuth != nil:
			p.authType = "basic"
		case len(p.Headers) > 0:
			p.authType = "headers"
		}
		p.Bearer = os.ExpandEnv(p.Bearer)
		for k, v := range p.Headers {
			p.Headers[k] = os.ExpandEnv(v)
		}
		if p.BasicAuth != nil {
			p.BasicAuth.Username = os.ExpandEnv(p.BasicAuth.Username)
			p.BasicAuth.Password = os.ExpandEnv(p.BasicAuth.Password)
		}
	}

	return &cfg, nil
}

func (c *Config) match(url string) *Profile {
	for i := range c.Profiles {
		if matchPattern(c.Profiles[i].Match, url) {
			return &c.Profiles[i]
		}
	}
	return nil
}

func matchPattern(pattern, url string) bool {
	if strings.HasSuffix(pattern, "**") {
		return strings.HasPrefix(url, strings.TrimSuffix(pattern, "**"))
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(url, strings.TrimSuffix(pattern, "*"))
	}
	return pattern == url
}

const exampleConfig = `# authcurl config
# Auth is auto-injected when the request URL matches a profile's pattern.
# Environment variables are supported: ${VAR_NAME}

profiles:
  # Bearer token auth
  - name: github
    match: "https://api.github.com/**"
    bearer: "${GITHUB_TOKEN}"

  # Custom headers
  - name: internal-api
    match: "https://api.example.com/**"
    headers:
      X-API-Key: "${API_KEY}"
      X-Custom: "some-value"

  # Basic auth
  - name: legacy
    match: "https://legacy.example.com/**"
    basic_auth:
      username: "${LEGACY_USER}"
      password: "${LEGACY_PASS}"
`

func initConfig() error {
	path := configPath()

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("config already exists at %s", path)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	if err := os.WriteFile(path, []byte(exampleConfig), 0o600); err != nil {
		return err
	}

	fmt.Printf("Config created at %s\n", path)
	fmt.Println("Edit it to add your profiles.")
	return nil
}
