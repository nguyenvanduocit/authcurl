package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		os.Exit(0)
	}

	if args[0] == "init" {
		if err := initConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "authcurl: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if args[0] == "list" {
		cfg, err := loadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "authcurl: %v\n", err)
			os.Exit(1)
		}
		listProfiles(cfg)
		return
	}

	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "authcurl: %v\n", err)
		fmt.Fprintln(os.Stderr, "Run 'authcurl init' to create a config file.")
		os.Exit(1)
	}

	url := extractURL(args)
	if url != "" {
		if profile := cfg.match(url); profile != nil {
			fmt.Fprintf(os.Stderr, "authcurl: using profile %q\n", profile.Name)
			args = injectAuth(args, profile)
		}
	}

	execCurl(args)
}

func printUsage() {
	fmt.Print(`authcurl - curl wrapper with automatic auth injection

Usage:
  authcurl [curl-options...] <url>    Execute curl with auto-injected auth
  authcurl init                       Create config file with examples
  authcurl list                       List configured profiles

Config: ~/.config/authcurl/config.yaml

Examples:
  authcurl https://api.github.com/user
  authcurl -X POST -d '{"title":"test"}' https://api.github.com/repos/owner/repo/issues
`)
}

func listProfiles(cfg *Config) {
	if len(cfg.Profiles) == 0 {
		fmt.Println("No profiles configured.")
		return
	}
	for _, p := range cfg.Profiles {
		fmt.Printf("  %-20s %-40s [%s]\n", p.Name, p.Match, p.authType)
	}
}
