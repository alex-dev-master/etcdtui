# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-12-30

### Added
- Interactive terminal UI for etcd
- Tree view for browsing etcd keys in hierarchical structure
- CRUD operations (Create, Read, Update, Delete) for keys
- Live watch mode to monitor key changes in real-time
- Prefix search functionality
- Multiple profile management for different etcd clusters
- Secure authentication support (username/password)
- TLS certificate support for secure connections
- Keyboard-driven navigation
- Profile selection screen
- Debug panel (F1) for troubleshooting
- Configuration file support (`~/.config/etcdtui/config.yaml`)
- Base64 password encoding in config
- Homebrew installation support
- Binary releases for macOS and Linux (amd64, arm64)

### Features
- **Profile Management**: Create, edit, delete, and switch between etcd cluster profiles
- **Key Operations**: Navigate, view, edit, and delete keys with syntax highlighting
- **Watch Mode**: Real-time monitoring of key changes with PUT/DELETE events
- **Search**: Quick prefix-based search for keys
- **Multi-platform**: Support for macOS and Linux on both AMD64 and ARM64 architectures

[0.1.0]: https://github.com/alex-dev-master/etcdtui/releases/tag/v0.1.0
