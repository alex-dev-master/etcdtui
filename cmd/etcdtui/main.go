package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alex-dev-master/etcdtui/internal/app/layouts"
	"github.com/spf13/pflag"
)

// Build-time variables (set via ldflags)
var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

var (
	profileName = pflag.StringP("profile", "p", "", "Profile name to use for connection")
	showHelp    = pflag.BoolP("help", "h", false, "Show help message")
	showVersion = pflag.BoolP("version", "v", false, "Show version")
)

func main() {
	pflag.Parse()

	if *showVersion {
		printVersion()
		return
	}

	if *showHelp {
		printUsage()
		return
	}

	m := layouts.NewManager()

	// Set profile from CLI flag
	if *profileName != "" {
		m.SetProfileName(*profileName)
	}

	ctx := context.Background()
	if err := m.Render(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printVersion() {
	fmt.Printf("etcdtui %s\n", version)
	fmt.Printf("  commit:  %s\n", commit)
	fmt.Printf("  built:   %s\n", buildTime)
}

func printUsage() {
	fmt.Printf("etcdtui %s - Interactive TUI for etcd\n", version)
	fmt.Print(`
Usage:
  etcdtui [flags]

Flags:
  -p, --profile string   Profile name to use for connection
  -h, --help             Show help message
  -v, --version          Show version

Config file: ~/.config/etcdtui/config.yaml

Example config:

profiles:
  - name: local
    endpoints: ["localhost:2379"]
    default: true

  - name: production
    endpoints: ["etcd1.prod:2379", "etcd2.prod:2379"]
    username: admin
    password: "base64:YWRtaW4xMjM="
    tls:
      enabled: true
      ca_file: "/path/to/ca.crt"
      cert_file: "/path/to/client.crt"
      key_file: "/path/to/client.key"
`)
}
