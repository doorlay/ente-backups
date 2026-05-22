# Contributing to ente-backups

## Development Setup

1. Clone the repository and navigate to it:
   ```bash
   git clone https://github.com/okwn/ente-backups.git
   cd ente-backups
   ```

2. Copy the environment template and configure your credentials:
   ```bash
   cp .env.example .env
   # Edit .env with your ENTE_* credentials and EXPORT_DIR path
   ```

3. Start the service with Docker Compose:
   ```bash
   docker compose up --build
   ```

## Manual Operations

### Manual Sync

To trigger a manual sync of your ente library:
```bash
docker compose exec backups ente-sync
```

### Manual Backup

To run a manual restic backup:
```bash
docker compose exec backups bash /usr/local/bin/restic-backup.sh
```

### Viewing Logs

Follow logs in real-time:
```bash
docker compose logs -f
```

## Testing

**Currently, there are no automated tests for this project.**

Before submitting a PR, ensure:
- `docker compose config` passes validation
- Manual sync and backup commands work as expected
- Any new functionality includes appropriate test coverage

## Project Structure

- `main.go` — Go service that monitors photos and triggers restic backups
- `restic-backup.sh` — Shell script that runs restic backup with notification
- `Dockerfile` — Container image definition
- `docker-compose.yml` — Service orchestration