# etcdtui

Interactive terminal UI for etcd3 - browse, edit, and monitor your etcd cluster with ease.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Features

- 🌲 **Tree View** - Browse etcd keys in a hierarchical tree structure
- 📝 **Edit Values** - Edit JSON/YAML values with syntax validation
- 👀 **Live Watch** - Monitor key changes in real-time
- 🔒 **Locks Dashboard** - View and manage distributed locks
- 🔐 **Secure Auth** - Support for username/password and TLS certificates
- 📋 **Multiple Profiles** - Quick switching between different etcd clusters
- ⌨️ **Keyboard-Driven** - Efficient navigation with vim-like shortcuts

## Installation

### From source

```bash
git clone https://github.com/alexandr/etcdtui.git
cd etcdtui
make build
```

### Using go install

```bash
go install github.com/alexandr/etcdtui/cmd/etcdtui@latest
```

## Quick Start

```bash
# Start with default connection (localhost:2379)
etcdtui

# Use specific profile
etcdtui --profile production

# Connect directly
etcdtui --endpoints etcd1:2379,etcd2:2379 --username admin
```

## Configuration

Create `~/.config/etcdtui/config.yaml`:

```yaml
profiles:
  production:
    endpoints:
      - etcd-prod1.example.com:2379
      - etcd-prod2.example.com:2379
    username: admin
    password: secret
    tls:
      enabled: true
      cert: /path/to/client.crt
      key: /path/to/client.key
      ca: /path/to/ca.crt

  staging:
    endpoints:
      - etcd-staging.example.com:2379
    username: admin
    password: secret

  local:
    endpoints:
      - localhost:2379
```

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `j/k` or `↓/↑` | Navigate tree |
| `Enter` | Expand/collapse node or edit value |
| `n` | New key |
| `d` | Delete key |
| `e` | Edit value |
| `w` | Watch mode |
| `l` | Locks dashboard |
| `r` | Refresh |
| `/` | Search |
| `?` | Show help |
| `q` or `Ctrl+C` | Quit |

## Screenshots

### Main View
```
┌─ Keys ───────────────┐ ┌─ Details ─────────────────────────────┐
│                      │ │ Key: /services/api/v1/config          │
│ ▼ /services          │ │                                       │
│   ▼ /api             │ │ {                                     │
│     • v1/config   ●  │ │   "port": 8080,                       │
│     • v1/endpoints   │ │   "timeout": 30,                      │
│   ▶ /auth            │ │   "debug": false                      │
│ ▼ /config            │ │ }                                     │
│   • database-url     │ │                                       │
│   • redis-url        │ │ Revision: 12345                       │
│ ▼ /locks             │ │ Modified: 2024-12-28 10:15:23         │
│   🔒 payment [28s]   │ │ TTL: ∞                                │
└──────────────────────┘ └───────────────────────────────────────┘
```

## Building from Source

```bash
# Clone the repository
git clone https://github.com/alexandr/etcdtui.git
cd etcdtui

# Build
make build

# Run
./bin/etcdtui
```

## Development

```bash
# Run without building
make run

# Run tests
make test

# Run linter
make lint

# Clean build artifacts
make clean
```

## Roadmap

### Milestone 1 - MVP ✅
- [x] Basic UI (tree + details)
- [ ] Connect to etcd
- [ ] Read keys
- [ ] Display values

### Milestone 2 - Core Features
- [ ] Edit values
- [ ] Delete keys
- [ ] Create new keys
- [ ] Connection profiles
- [ ] TLS support

### Milestone 3 - Advanced Features
- [ ] Watch mode
- [ ] Locks dashboard
- [ ] Search functionality
- [ ] Revision history
- [ ] Bulk operations
- [ ] Export/Import (YAML/JSON)

## Similar Projects

- [k9s](https://k9scli.io/) - Kubernetes CLI
- [lazygit](https://github.com/jesseduffield/lazygit) - Git TUI
- [lazydocker](https://github.com/jesseduffield/lazydocker) - Docker TUI

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details

## Author

Alexandr - [@alex-dev-master](https://github.com/alex-dev-master)
