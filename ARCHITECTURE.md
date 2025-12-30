# Architecture

This document describes the project structure and architectural decisions.

## Directory Structure

```
etcdtui/
├── cmd/
│   └── etcdtui/
│       └── main.go                 # Application entry point, CLI flags
│
├── internal/
│   ├── app/
│   │   ├── actions/                # Business logic and action handlers
│   │   │   ├── general/
│   │   │   │   ├── state.go        # State management for main view
│   │   │   │   ├── etcd.go         # etcd operations (CRUD, refresh)
│   │   │   │   └── actions.go      # User action handlers (edit, delete)
│   │   │   │
│   │   │   └── profiles/
│   │   │       ├── state.go        # State management for profiles view
│   │   │       └── actions.go      # Profile form handlers
│   │   │
│   │   ├── layouts/                # UI layout and input routing
│   │   │   ├── manager.go          # Layout manager, screen switching
│   │   │   ├── general/
│   │   │   │   └── layout.go       # Main view layout
│   │   │   └── profiles/
│   │   │       └── layout.go       # Profile selection layout
│   │   │
│   │   └── connection/             # Connection management
│   │       └── etcd/
│   │           └── manager.go      # etcd connection manager
│   │
│   ├── config/                     # Configuration management
│   │   ├── config.go               # Config loading/saving with Viper
│   │   ├── profile.go              # Profile struct and encoding
│   │   └── errors.go               # Config errors
│   │
│   └── ui/
│       └── panels/                 # Reusable UI components
│           ├── keys/               # Keys tree panel
│           ├── details/            # Key details panel
│           ├── statusbar/          # Status bar panel
│           └── debug/              # Debug log panel
│
└── pkg/
    └── etcd/
        ├── client.go               # etcd client wrapper
        ├── config.go               # Client configuration
        ├── kv.go                   # Key-value operations
        ├── watch.go                # Watch operations
        └── ...                     # Other etcd operations
```

## Application Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                         main.go                                      │
│  - Parse CLI flags (--profile, --help)                              │
│  - Create layout manager                                             │
│  - Start application                                                 │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     layouts/manager.go                               │
│                                                                      │
│  - Load configuration                                                │
│  - If --profile flag: go directly to general layout                 │
│  - Otherwise: show profiles layout first                            │
│  - Handle screen switching                                           │
└─────────────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┴───────────────┐
              ▼                               ▼
┌─────────────────────────────┐  ┌─────────────────────────────┐
│   layouts/profiles/          │  │   layouts/general/           │
│                              │  │                              │
│  - Profile list              │  │  - Keys tree                 │
│  - Profile details           │  │  - Key details               │
│  - Create/Edit/Delete        │  │  - Status bar                │
│                              │  │                              │
│  On connect ─────────────────┼──│                              │
│                              │  │  On 'p' key ─────────────────│
└─────────────────────────────┘  └─────────────────────────────┘
```

## Package Responsibilities

### `cmd/etcdtui/`

Application entry point:
- Parse CLI flags (`--profile`, `--help`, `--version`)
- Create and run layout manager

### `internal/config/`

Configuration management:
- Load/save config from `~/.config/etcdtui/config.yaml`
- Profile struct with endpoints, auth, TLS settings
- Password encoding (base64)

### `internal/app/actions/`

**What happens** — Business logic and action handlers.

| Package | File | Responsibility |
|---------|------|----------------|
| `general` | `state.go` | Main view state: panels, connection, current key |
| `general` | `etcd.go` | etcd operations: connect, list, CRUD, refresh |
| `general` | `actions.go` | User actions: edit form, delete modal, search |
| `profiles` | `state.go` | Profiles view state: selected profile, UI components |
| `profiles` | `actions.go` | Profile actions: create/edit form, delete modal |

### `internal/app/layouts/`

**How it's arranged** — UI layout composition and input routing.

| File | Responsibility |
|------|----------------|
| `manager.go` | Creates layouts, handles screen switching |
| `general/layout.go` | Main view: flex layout, keyboard routing |
| `profiles/layout.go` | Profile selection: list, details, keyboard routing |

### `internal/ui/panels/`

**How it looks** — Reusable UI components.

Each panel is a self-contained UI component with its own:
- Visual representation (TextView, TreeView, Form, etc.)
- Local state
- Public API for interaction

### `pkg/etcd/`

**External interface** — etcd client wrapper.

Low-level etcd operations:
- Connection handling with TLS support
- CRUD operations
- Watch functionality
- Lease management
- Cluster status

## Design Principles

1. **Separation of Concerns**
   - Layouts handle visual composition and input routing
   - Actions handle business logic and user interactions
   - Panels are reusable UI components
   - Config handles persistence
   - pkg/etcd handles external communication

2. **Consistent Structure**
   - Both `general` and `profiles` follow the same pattern
   - `state.go` for state management
   - `actions.go` for action handlers
   - `layout.go` for UI composition

3. **Dependency Direction**
   ```
   cmd/etcdtui
        │
        ▼
   layouts/manager
        │
   ┌────┴────┐
   ▼         ▼
 layouts   layouts
 /general  /profiles
   │         │
   ▼         ▼
 actions   actions
 /general  /profiles
   │         │
   ▼         ▼
 ui/panels   config
   │
   ▼
 pkg/etcd
   ```

4. **Testability**
   - State can be mocked for testing actions
   - Actions can be tested independently of UI
   - Panels can be tested in isolation
   - Config can be tested with temp files
