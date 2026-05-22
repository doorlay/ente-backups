# 06_SELECTED_5_PR_PLAN.md — ente-backups

## Selected PRs (in priority order)

### PR 1: Add healthcheck to docker-compose.yml (Reliability)
**Files**: `docker-compose.yml`
**Change**: Add `healthcheck` block under `backups` service:
```yaml
healthcheck:
  test: ["CMD", "test", "-f", "/srv/backups/ente-sync.lock"]
  interval: 5m
  timeout: 10s
  retries: 3
```
**Why**: Detect silent binary failures; current `restart: unless-stopped` only reacts to exit codes, not hangs
**Risk**: LOW — harmless check, only affects container restart behavior
**CI needed**: Validate docker-compose config

### PR 2: Add CONTRIBUTING.md (Onboarding)
**Files**: `CONTRIBUTING.md` (new)
**Change**: Add a contributing guide covering:
- Dev setup (clone, docker compose up)
- Manual sync: `docker compose exec backups ente-sync`
- Manual backup: `docker compose exec backups bash /usr/local/bin/restic-backup.sh`
- View logs: `docker compose logs -f`
- No tests currently; add tests before PRs
**Why**: Helps new contributors understand project structure and constraints
**Risk**: NONE — docs only
**CI needed**: None (or markdown linter)

### PR 3: Add Go build CI (Quality Gate)
**Files**: `.github/workflows/go.yml` (new)
**Change**: GitHub Actions workflow:
```yaml
name: Build
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.23' }
      - run: go build -o ente-sync .
```
**Why**: Catch build failures before merge; enforce `go.mod tidy` clean
**Risk**: NONE — only validates build, no code changes
**CI needed**: GitHub Actions (free for public repos)

### PR 4: Improve notify() HTTP error handling (main.go)
**Files**: `main.go`
**Change**: Check HTTP status code in `notify()`:
```go
resp, err := http.Post(url, "text/plain", strings.NewReader(msg))
if err != nil {
    log.Printf("Failed to send notification: %v", err)
    return
}
defer resp.Body.Close()
if resp.StatusCode >= 400 {
    log.Printf("Notification failed with status %d", resp.StatusCode)
}
```
**Why**: Silent failures in notification code make debugging harder; non-200 responses should be logged
**Risk**: LOW — notification utility only, no photo/restore/encryption code touched
**CI needed**: Go build validation

## PR Execution Order
1. PRs #1, #2, #3 can be done in any order (no file conflicts)
2. PR #4 last (small Go change, touches main.go notification function only)

## What NOT to touch
- `main.go`: sync/export, file lock, timeout, result batching, encryption
- `restic-backup.sh`: backup/restore/prune/encryption
- `Dockerfile`: base images, versions, apt packages
- `crontab`: schedule
- `LICENSE`: MIT (do not change)