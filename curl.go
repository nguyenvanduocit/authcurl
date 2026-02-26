package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func extractURL(args []string) string {
	for _, arg := range args {
		if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
			return arg
		}
		if strings.HasPrefix(arg, "--url=") {
			return strings.TrimPrefix(arg, "--url=")
		}
	}
	return ""
}

func injectAuth(args []string, profile *Profile) []string {
	var extra []string

	if profile.Bearer != "" {
		extra = append(extra, "-H", "Authorization: Bearer "+profile.Bearer)
	}

	for key, value := range profile.Headers {
		extra = append(extra, "-H", key+": "+value)
	}

	if profile.BasicAuth != nil {
		extra = append(extra, "-u", profile.BasicAuth.Username+":"+profile.BasicAuth.Password)
	}

	return append(extra, args...)
}

func execCurl(args []string) {
	curlPath, err := exec.LookPath("curl")
	if err != nil {
		fmt.Fprintln(os.Stderr, "authcurl: curl not found in PATH")
		os.Exit(1)
	}

	// syscall.Exec replaces the current process — exit codes, signals,
	// stdin/stdout/stderr all pass through to curl directly.
	if err := syscall.Exec(curlPath, append([]string{"curl"}, args...), os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "authcurl: exec failed: %v\n", err)
		os.Exit(1)
	}
}
