#!/usr/bin/env bash
set -euo pipefail

BACKUP_PATH="/data/ente-photos"

if [[ -z "${RESTIC_REPOSITORY:-}" ]]; then
  echo "RESTIC_REPOSITORY not set, skipping backup"
  exit 0
fi

notify() {
  if [[ -n "${NTFY_TOPIC:-}" ]]; then
    curl -s -d "$1" "https://ntfy.sh/${NTFY_TOPIC}" >/dev/null
  fi
}

echo "Starting restic backup of ${BACKUP_PATH}"

if ! restic backup "${BACKUP_PATH}"; then
  notify "Restic backup failed"
  exit 1
fi

echo "Running restic forget/prune"

if ! restic forget --prune --keep-daily 30 --keep-weekly 26 --keep-monthly 12; then
  notify "Restic prune failed"
  exit 1
fi

echo "Restic backup completed successfully"
notify "Restic backup completed successfully"
