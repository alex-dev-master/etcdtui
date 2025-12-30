package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alexandr/etcdtui/internal/app/layouts"
	"github.com/spf13/pflag"
)

var (
	profileName = pflag.StringP("profile", "p", "", "Profile name to use for connection")
	showHelp    = pflag.BoolP("help", "h", false, "Show help message")
)

// Version is set during build

func main() {
	pflag.Parse()

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
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`etcdtui - Interactive TUI for etcd

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
