# 00_STATE.md — ente-backups

## Repository Identity
- **Name**: ente-backups
- **Owner**: doorlay (upstream) / okwn (fork)
- **URL**: https://github.com/doorlay/ente-backups
- **Description**: Automated local and remote backups of your Ente photo library
- **Language**: Go
- **License**: MIT
- **Archived**: No
- **Default branch**: main

## Fork Status
- Forked: Yes (2026-05-22)
- Fork URL: https://github.com/okwn/ente-backups
- Upstream remote: `upstream` → https://github.com/doorlay/ente-backups
- Fork is ahead of upstream: No (identical HEAD)

## Repository Structure
```
ente-backups/
├── main.go           # Go binary: ente-sync (hourly Ente export + ntfy notifications)
├── restic-backup.sh  # Bash script: daily S3 backup via restic
├── crontab           # Supercronic schedule (hourly sync, daily backup)
├── Dockerfile        # Multi-stage build: Go app + debian-slim + ente-cli + restic
├── docker-compose.yml # Docker Compose config (no healthchecks)
├── .env.example      # Environment template
├── go.mod            # Go module (1.23)
├── .gitignore        # Ignores *.env only
├── README.md         # Setup guide (74 lines)
└── LICENSE           # MIT License
```

## Tech Stack
- **Language**: Go 1.23
- **Container**: Debian Bookworm slim + Docker Compose
- **Backup Tools**: Ente CLI (ente), Restic
- **Scheduler**: Supercronic (cron in containers)
- **Notifications**: ntfy.sh (HTTP POST)

## Key Observations
- Small, single-binary project with no external Go dependencies
- Only one Go file (main.go) with ~138 lines
- No tests, no CI/CD, no linting
- No CONTRIBUTING.md
- Upstream .gitignore includes CLAUDE.md, fork's does not
- Hardcoded version numbers in Dockerfile (supercronic v0.2.33, ente-cli v0.2.3)
- No healthchecks in docker-compose.yml

## CI/CD: None
## Test Framework: None
## Build Commands: `go build -o ente-sync .`
## Lint: None

## Risk Areas (DO NOT MODIFY)
- main.go: `notify()` — HTTP notification logic, recordResult
- restic-backup.sh: `restic backup/forget/prune` — encryption, backup, restore logic
- main.go: `cmd.Run()` with `ente export` — photo data handling

## External Services
- Ente API (ente CLI)
- S3-compatible storage (restic)
- ntfy.sh (notifications)