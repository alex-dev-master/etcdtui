# Architecture

This document describes the project structure and architectural decisions.

## Directory Structure

```
etcdtui/
├── cmd/
│   └── etcdtui/
│       └── main.go              # Application entry point
│
├── internal/
│   ├── app/
│   │   ├── actions/             # Business logic and action handlers
│   │   │   └── general/
│   │   │       ├── state.go     # State management (panels, connection, flags)
│   │   │       ├── etcd.go      # etcd operations (CRUD, refresh, etc.)
│   │   │       └── actions.go   # User action handlers (edit, delete, etc.)
│   │   │
│   │   ├── layouts/             # UI layout and input routing
│   │   │   ├── manager.go       # Layout manager
│   │   │   └── general/
│   │   │       └── layout.go    # Flex layouts, input capture routing
│   │   │
│   │   └── connection/          # Connection management
│   │       └── etcd/
│   │           └── manager.go   # etcd connection manager
│   │
│   └── ui/
│       └── panels/              # Reusable UI components
│           ├── keys/            # Keys tree panel
│           ├── details/         # Key details panel
│           ├── statusbar/       # Status bar panel
│           └── debug/           # Debug log panel
│
└── pkg/
    └── etcd/
        └── client.go            # etcd client wrapper
```

## Package Responsibilities

### `internal/app/actions/`

**What happens** — Business logic and action handlers.

| File | Responsibility |
|------|----------------|
| `state.go` | Holds application state: panels, connection manager, current key, UI flags |
| `etcd.go` | etcd operations: connect, list, get, put, delete, refresh |
| `actions.go` | User action handlers: edit form, delete modal, help dialog |

### `internal/app/layouts/`

**How it's arranged** — UI layout composition and input routing.

| File | Responsibility |
|------|----------------|
| `manager.go` | Creates and manages layouts |
| `layout.go` | Flex composition, keyboard input routing to actions |

### `internal/ui/panels/`

**How it looks** — Reusable UI components.

Each panel is a self-contained UI component with its own:
- Visual representation (TextView, TreeView, Form, etc.)
- Local state
- Public API for interaction

### `pkg/etcd/`

**External interface** — etcd client wrapper.

Low-level etcd operations, connection handling, data types.

## Data Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                        User Input                                │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    layouts/general/layout.go                     │
│                                                                  │
│  - Receives keyboard events via SetInputCapture                  │
│  - Routes events to appropriate action handlers                  │
│  - Manages Flex layout composition                               │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    actions/general/                              │
│                                                                  │
│  state.go   │ Holds panels, connection, current state           │
│  actions.go │ Handles user actions (edit, delete, etc.)         │
│  etcd.go    │ Performs etcd operations                          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    ui/panels/                                    │
│                                                                  │
│  keys/      │ Displays key tree                                 │
│  details/   │ Shows key details and action buttons              │
│  statusbar/ │ Shows status and hints                            │
│  debug/     │ Shows debug logs                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    pkg/etcd/                                     │
│                                                                  │
│  - Low-level etcd client operations                              │
│  - Connection management                                         │
└─────────────────────────────────────────────────────────────────┘
```

## Design Principles

1. **Separation of Concerns**
   - Layouts handle visual composition and input routing
   - Actions handle business logic and user interactions
   - Panels are reusable UI components
   - pkg/etcd handles external communication

2. **Single Responsibility**
   - Each file has a clear, focused purpose
   - State management is centralized in `state.go`
   - etcd operations are isolated in `etcd.go`

3. **Dependency Direction**
   - `layouts` → `actions` → `panels` → `pkg`
   - Higher layers depend on lower layers, not vice versa

4. **Testability**
   - State can be mocked for testing actions
   - Actions can be tested independently of UI
   - Panels can be tested in isolation
