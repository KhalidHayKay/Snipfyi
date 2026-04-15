# Smply — URL Shortener

A production-style URL shortener built in Go, featuring a public API, email-based API key authentication, rate limiting, and background job support.

---

## What It Does

- Shortens URLs via Base62 encoding
- Exposes a public REST API (`POST /api/shorten`, `GET /:code`)
- Issues API keys through an email + magic link flow
- Enforces rate limiting and abuse prevention

### API Key Flow

1. Request access via email on the `/api` page
2. Receive a magic link (valid for 15 minutes, single-use)
3. Click the link to activate and receive your API key (shown once — store it)
4. Pass the key as an `Authorization` header on API requests

---

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **Cache / Rate limiting:** Redis
- **Container:** Docker + Docker Compose

---

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- Go 1.25+ (only if running outside Docker)

---

## Local Development Setup

### 1. Clone the repo

```bash
git clone https://github.com/your-org/smply.git
cd smply
```

### 2. Set up your environment file

```bash
cp .env.example .env
```

Open `.env` and fill in the required values (SMTP credentials, app secret, etc.). Make sure `APP_ENV=development` is set.

### 3. Create the dev Compose override file

The repo ships a base `compose.yaml` that is production-safe (no volume mounts). For local development, you need a `compose.override.yaml` that mounts your source code and enables hot reload via [Air](https://github.com/air-verse/air).

Create `compose.override.yaml` in the project root:

```yaml
services:
 app:
  volumes:
   - ./:/app
   - go-module-cache:/root/go/pkg/mod
  worker:
    volumes:
   - ./:/app
   - go-module-cache:/root/go/pkg/mod

volumes:
 go-module-cache:
```

> **Why is this not committed?**
> The override mounts your local source code into the container, which is only meaningful in development. Keeping it out of the repo ensures the production server can run `docker compose up` safely without any extra flags.

### 4. Build and start

```bash
docker compose up --build
```

Docker Compose automatically merges `compose.override.yaml` with `compose.yaml`, so you don't need any extra flags. The app will start with Air watching for file changes and hot-reloading on save.

The app will be available at: `http://localhost:8000`

---

## Services

| Service | Port   | Description                 |
| ------- | ------ | --------------------------- |
| `app`   | `8000` | Go application              |
| `db`    | `5432` | PostgreSQL 16               |
| `redis` | `6379` | Redis 8 (password: `admin`) |

---

## Production Deployment

On the production server, only `compose.yaml` is present (no override file). Start the app with:

```bash
APP_ENV=production docker compose up --build -d
```

The Dockerfile detects `APP_ENV` at build time: in production it compiles the Go binary; in development it installs Air for hot reload. The startup script (`start.sh`) then runs the appropriate process at container start.

---

## API Reference

### Shorten a URL

```
POST /api/shorten
X-API-Key: <your-api-key>
Content-Type: application/json

{
  "url": "https://example.com/very/long/url"
}
```

**Response:**

```json
{
	"short_url": "https://smply.app/aB3xZ"
}
```

### Redirect

```
GET /:code
```

Redirects to the original URL.

---

## Security Notes

- API keys and magic tokens are stored hashed (SHA-256) — raw values are never persisted
- Magic tokens expire after 15 minutes and are single-use
- Only one active API key per email; only one active token per email
- API keys expire after 30 days
