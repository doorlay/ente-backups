### Overview
Automated local and remote backups of your Ente photo library.
- Hourly syncing via Ente CLI
- Point-in-time encrypted backups in S3 via restic (optional)
- Push notifications for backup success/errors via ntfy.sh (optional)

### Prerequisites
- Docker Compose 

### Setup
This is best run on a headless server with enough free storage for your photo library (e.g. I run this on a Raspberry Pi connected to a 240GB SSD). By default, syncing from Ente occurs hourly and backups to S3 happen daily at 8pm PST, but these are both configurable.

1. Copy the example environment file and fill in the Ente section:
```
cp .env.example .env
```
If you want photos exported somewhere other than `/data/ente-photos`, set `EXPORT_DIR` in `.env` to your chosen path; make sure the directory exists on the host first (`sudo mkdir -p "$EXPORT_DIR"`).

2. Build and start the container:
```
docker compose up -d --build
```
3. Perform the initial Ente login:
```
docker compose exec backups ente account add
```
You'll be prompted for your Ente email, password, and an export directory. **For the export directory, enter the same value as `EXPORT_DIR` in `.env`** (default `/data/ente-photos`); this must be a path that exists *inside* the container, which is what the bind mount in step 1 sets up. Anything else will fail with `invalid export directory`.

4. Kick off the first sync to verify (otherwise it'll wait until the top of the next hour):
```
docker compose exec backups ente-sync
```
The first full sync can take hours depending on library size; you can detach the terminal with Ctrl+P, Ctrl+Q without stopping the export.

### Setting up Restic (optional)
Restic creates encrypted, deduplicated backups of your photo library to an S3 bucket. 

1. Create an S3 bucket in AWS (or any S3-compatible provider like Backblaze B2 or Wasabi).
2. Create an IAM user with read/write access to the bucket and generate an access key.
3. Fill in the Restic section of `.env`:
```
RESTIC_REPOSITORY=s3:s3.amazonaws.com/your-bucket
RESTIC_PASSWORD=your-encryption-password    # used to encrypt backups; store this safely, e.g. in a password manager
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```
4. Initialize the restic repository:
```
docker compose exec backups restic init
```

If the Restic variables are left empty, the daily backup job will skip silently.

### Setting up ntfy.sh notifications (optional)
[ntfy.sh](https://ntfy.sh) sends push notifications to your phone or desktop when syncs or backups succeed or fail.

1. Install the ntfy app wherever you plan to receive notifications. 
2. Choose a topic name; this can be anything, but should be unique and hard to guess (e.g. `my-ente-backups-a1b2c3`).
3. Subscribe to your topic name in the ntfy app.
4. Set `NTFY_TOPIC` in `.env` to your topic name:
```
NTFY_TOPIC=my-ente-backups-a1b2c3
```

If `NTFY_TOPIC` is left empty, notifications are disabled.

Sync failures and timeouts trigger an immediate notification. To prevent notification fatigue, successes are batched into a daily summary sent after every 24 runs, reporting how many succeeded. Restic backups notify on both success and failure.

### Development
- `docker compose exec backups ente-sync`: run the ente export manually
- `docker compose exec backups bash /usr/local/bin/restic-backup.sh`: run the restic backup manually

### Notes
- Logs: `docker compose logs -f`

## Contributing

Contributions are welcome!
