### Overview
Automated local and S3 backups of your Ente photo library.
- Hourly syncing via Ente CLI
- Point-in-time encrypted backups in S3 via restic (optional)
- Push notifications for backup success/errors via ntfy.sh (optional)

### Setup
This is best run on a headless server with enough free storage for your photo library (e.g. I run this on a Raspberry Pi connected to a 240GB SSD). By default, syncing from Ente occurs hourly, and backups to S3 happen daily at 8pm PST.

1. Copy the example environment file and fill in the Ente section:
```
cp .env.example .env
```
2. Build and start the container:
```
docker compose up -d --build
```
3. Perform the initial Ente login:
```
docker compose exec backups ente account add
```

### Setting up Restic (optional)
Restic creates encrypted, deduplicated backups of your photo library to an S3 bucket.

1. Create an S3 bucket in AWS (or any S3-compatible provider like Backblaze B2 or Wasabi).
2. Create an IAM user with read/write access to the bucket and generate an access key.
3. Fill in the Restic section of `.env`:
```
RESTIC_REPOSITORY=s3:s3.amazonaws.com/your-bucket
RESTIC_PASSWORD=your-encryption-password    # used to encrypt backups; store this safely
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```
4. Initialize the restic repository (one-time setup):
```
docker compose exec backups restic init
```

If the Restic variables are left empty, the daily backup job will skip silently.

### Setting up ntfy.sh notifications (optional)
[ntfy.sh](https://ntfy.sh) sends push notifications to your phone or desktop when syncs or backups succeed or fail.

1. Install the ntfy app on your phone ([Android](https://play.google.com/store/apps/details?id=io.heckel.ntfy) / [iOS](https://apps.apple.com/app/ntfy/id1625396347)) or use the [web app](https://ntfy.sh/app).
2. Choose a topic name; this can be anything, but should be unique and hard to guess (e.g. `my-ente-backups-a1b2c3`).
3. Subscribe to your topic in the ntfy app.
4. Set `NTFY_TOPIC` in `.env`:
```
NTFY_TOPIC=my-ente-backups-a1b2c3
```

If `NTFY_TOPIC` is left empty, notifications are disabled.

Sync failures and timeouts trigger an immediate notification. To prevent notification fatigue, successes are batched into a daily summary sent after every 24 runs, reporting how many succeeded. Restic backups notify on both success and failure.

### Development
- `docker compose exec backups ente-sync` — run the ente export manually
- `docker compose exec backups bash /usr/local/bin/restic-backup.sh` — run the restic backup manually

### Notes
- Logs: `docker compose logs -f`
