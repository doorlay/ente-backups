# 01_REPO_MAP.md — ente-backups

## What It Does
Automated local and remote backups of an Ente photo library:
1. **Hourly sync**: Ente CLI exports photos to a local directory (`EXPORT_DIR`)
2. **Daily backup**: Restic backs up the local directory to S3 with encryption
3. **Notifications**: ntfy.sh push notifications for sync/backup success or failure

## Tech Stack
- **Go 1.23** — single binary (`ente-sync`)
- **Bash** — restic backup script
- **Docker** — multi-stage Dockerfile, docker-compose
- **Supercronic** — cron scheduling inside container
- **External tools**: Ente CLI (ente), Restic, curl

## Directory Structure
```
ente-backups/
├── main.go           # ente-sync Go binary (lock file, env validation, timeout, notifications)
├── restic-backup.sh  # Daily S3 backup script (backup + forget/prune)
├── crontab           # Supercronic schedule
│                        0 * * * * → hourly ente-sync
│                        0 20 * * * → daily restic-backup.sh
├── Dockerfile        # Multi-stage: golang:1.23-bookworm builder → debian:bookworm-slim runtime
├── docker-compose.yml # Service: backups (no healthchecks, restart: unless-stopped)
├── .env.example      # 4 sections: Ente sync, Restic S3, ntfy, optional ALBUMS/INCLUDE_HIDDEN
├── go.mod            # module ente-backups, go 1.23
├── .gitignore        # *.env only
├── README.md         # Setup instructions, development commands
└── LICENSE           # MIT
```

## Runtime Behavior
- `ente-sync` runs hourly via cron, uses file lock to prevent concurrent runs
- 6-hour timeout on Ente export; on timeout → immediate ntfy notification
- Success results written to `/srv/backups/ente-sync-results.log`
- After 24 successful runs (or if failures exist), daily batch notification sent, log truncated
- `restic-backup.sh` runs daily; skips entirely if `RESTIC_REPOSITORY` is empty

## Build/Lint/Test Commands
```bash
# Build
go build -o ente-sync .

# Lint: None
# Test: None

# Dev commands (from README)
docker compose exec backups ente-sync         # manual sync
docker compose exec backups bash /usr/local/bin/restic-backup.sh  # manual backup
docker compose logs -f                        # view logs
```

## Critical Behavior Notes
- `main.go`: EXPORT_DIR must be set or binary exits with fatal
- `main.go`: Ente export path inside container MUST match EXPORT_DIR env var (bind mount)
- `main.go`: ntfy notifications only fire if NTFY_TOPIC is set; successes batched every 24h
- `restic-backup.sh`: silently skips if RESTIC_REPOSITORY is unset
- `restic-backup.sh`: always runs forget/prune after backup (no dry-run option)

## Files NOT to Modify (Risk Areas)
- `main.go` — photo sync logic, file lock, timeout, notification batching, encryption handling
- `restic-backup.sh` — backup/restore/encryption logic
- Dockerfile — base images, apt packages, tool versions
- docker-compose.yml — volume mounts, restart policy

## Quality Issues
- No tests
- No CI/CD
- No linting
- No CONTRIBUTING.md
- Fork's .gitignore missing `CLAUDE.md` (present in upstream)
- Hardcoded version numbers in Dockerfile (no update mechanism)
- docker-compose.yml lacks healthchecks
- No validation of S3 credentials or Ente auth before cron runs