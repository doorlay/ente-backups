# 05_PR_CANDIDATES.md — ente-backups

## Open Issues / PRs
- **None** — repo has 0 open issues and 0 open PRs

## Code Quality Issues (Safe to Fix)

### 1. main.go: HTTP response body not checked before closing
**File**: main.go, line 133
**Issue**: `http.Post` response is not checked for HTTP errors; response body is read-only closed with `resp.Body.Close()`, no status code check
**Risk**: LOW (notification utility only)
**Fix**: Check `resp.StatusCode` before closing; log if non-200

### 2. restic-backup.sh: curl output not checked
**File**: restic-backup.sh, line 13
**Issue**: `curl -s` result discarded; failure of curl notification silently ignored
**Risk**: LOW (notification only)
**Fix**: Check curl exit code or redirect stderr

### 3. .gitignore — already in sync with upstream
**File**: .gitignore
**Status**: Already matches upstream (contains both `*.env` and `CLAUDE.md`)

### 4. Dockerfile: hardcoded tool versions with no update mechanism
**File**: Dockerfile, lines 17, 22
**Issue**: supercronic v0.2.33 and ente-cli v0.2.3 hardcoded; no renovate/dependabot
**Risk**: MEDIUM (tool becomes outdated, security issues unpatched)
**Fix**: Add comments or a check script for version updates

### 5. docker-compose.yml: no healthchecks
**File**: docker-compose.yml
**Issue**: No healthcheck defined; restart: unless-stopped won't detect a crashed binary
**Risk**: MEDIUM (silent failures)
**Fix**: Add healthcheck to docker-compose.yml

### 6. No CONTRIBUTING.md
**Risk**: LOW (project is simple)
**Fix**: Add CONTRIBUTING.md with setup/testing instructions

### 7. No CI/CD
**Risk**: LOW (manual process documented)
**Fix**: Add GitHub Actions workflow for build validation

### 8. No test infrastructure
**Risk**: MEDIUM (no regression detection)
**Fix**: Add Go tests for main.go (env validation, lock file, timeout logic)

### 9. Ente CLI version hardcoded in Dockerfile
**File**: Dockerfile, line 22
**Issue**: `ENTE_CLI_VERSION=v0.2.3` pinned with no update mechanism
**Risk**: MEDIUM (outdated Ente CLI)
**Fix**: Document how to update; consider checking for new releases

## Low-Risk Change Candidates (Priority Order)
1. Add healthcheck to docker-compose.yml
2. Add CONTRIBUTING.md
3. Add GitHub Actions CI for `go build` validation
4. Fix HTTP response status check in main.go notify()
5. restic-backup.sh curl: already has `set -euo pipefail` (sufficient)
6. Dockerfile: hardcoded tool versions — consider documenting update process