# Monorepo Restructure Plan

## Goal

Restructure cli-portainer into cli-tools monorepo supporting multiple tools with shared code.

## Target Structure

```
cli-tools/
├── build.sh
├── README.md
├── CLAUDE.md
├── .gitignore
│
├── go/
│   ├── go.work                 # Go workspace
│   ├── common/                 # Shared code (stub for now)
│   │   ├── go.mod
│   │   └── placeholder.go
│   └── portainer/              # portainer-cli
│       ├── go.mod
│       ├── cmd/portainer-cli/
│       └── pkg/portainer/
│
├── rust/
│   ├── Cargo.toml              # Workspace
│   ├── common/                 # Shared code (stub for now)
│   │   ├── Cargo.toml
│   │   └── src/lib.rs
│   └── portainer/              # portainer-cli
│       ├── cli/
│       └── lib/
│
└── docs/
    └── plans/
```

## Tasks

1. Create new directory structure
2. Move go/ contents to go/portainer/
3. Move rust/ contents to rust/portainer/
4. Create go/common/ stub
5. Create rust/common/ stub
6. Create go/go.work workspace file
7. Update rust/Cargo.toml workspace
8. Update go.mod paths
9. Update build.sh for new paths
10. Update README.md
11. Verify builds and tests pass

## Notes

- Root folder stays as `cli-portainer` locally (will be `cli-tools` on GitHub)
- Common code extraction deferred until second tool needs it
- Each tool can be Go-only, Rust-only, or both
