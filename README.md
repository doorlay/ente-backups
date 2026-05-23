# Ente Backups

Automated local and remote backups of your Ente photo library.

## 🚀 Quick Start

### Prerequisites
- Docker Compose (v2.17.0 or later recommended)
- Docker Engine 24.0+

### 1. Clone and configure

```bash
git clone https://github.com/doorlay/ente-backups.git
cd ente-backups
cp .env.example .env
```

### 2. Set your Ente credentials in `.env`

```env
ENTE_EMAIL=your@email.com
ENTE_PASSWORD=your-password
EXPORT_DIR=/data/ente-photos   # change if you want photos elsewhere
```

Make sure the export directory exists on the host:
```bash
sudo mkdir -p "/data/ente-photos"
```

### 3. Start the container

```bash
docker compose up -d --build
```

### 4. Authenticate with Ente

```bash
docker compose exec backups ente account add
```

Enter your email and password when prompted. For the export directory, enter `/data/ente-photos` (or whatever you set as `EXPORT_DIR`).

### 5. Run your first sync

```bash
docker compose exec backups ente-sync
```

The first sync can take hours for large libraries. You can safely detach with `Ctrl+P, Ctrl+Q`.

---

## Features

- Hourly syncing via Ente CLI
- Point-in-time encrypted backups in S3 via restic (optional)
- Push notifications for backup success/errors via ntfy.sh (optional)

## Optional: Set up encrypted S3 backups

Restic creates encrypted, deduplicated backups to an S3 bucket.

1. Create an S3 bucket in AWS (or any S3-compatible provider like Backblaze B2 or Wasabi).
2. Create an IAM user with read/write access to the bucket and generate an access key.
3. Fill in the Restic section of `.env`:
```env
RESTIC_REPOSITORY=s3:s3.amazonaws.com/your-bucket
RESTIC_PASSWORD=your-encryption-password
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```
4. Initialize the repository:
```bash
docker compose exec backups restic init
```

## Optional: Enable notifications

[ntfy.sh](https://ntfy.sh) sends push notifications when syncs or backups succeed or fail.

1. Install the ntfy app on your phone or desktop.
2. Choose a topic name (e.g. `my-ente-backups-a1b2c3`) and subscribe to it in the app.
3. Add to `.env`:
```env
NTFY_TOPIC=my-ente-backups-a1b2c3
```

Sync failures trigger an immediate notification. Successes are batched into a daily summary.

## Development

```bash
docker compose exec backups ente-sync    # run sync manually
docker compose exec backups bash /usr/local/bin/restic-backup.sh  # run backup manually
docker compose logs -f                   # watch logs
```