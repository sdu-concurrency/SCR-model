# dk-if-open

Open-source deployment of the **Supply Chain Resilience** assessment tool, combining the Go/PocketBase backend and Vue 3 frontend into a single Docker image.

---

## Architecture

```
dk-if-open/
├── back/          # Go backend (PocketBase)
├── web/           # Vue 3 frontend
├── Dockerfile     # Multi-stage build: Go backend + Vue 3 frontend
└── .env.example   # Template for build-time environment variables
```

The container runs PocketBase as both the REST API server and static file server, serving the compiled frontend from `/pb_public`.

---

## Prerequisites

- **Docker** ≥ 24
- **Git** ≥ 2.13

---

## Quick Start

### 1. Clone

```bash
git clone https://github.com/sdu-concurrency/SCR-model.git
cd SCR-model
```

### 2. Build

```bash
docker build -t dk-if-open .
```

### 3. Run

```bash
docker run -d \
  -p 8080:8080 \
  -v dk-if-open-data:/pb_data \
  --name dk-if-open \
  dk-if-open
```

The application is available at `http://localhost:8080/`.

PocketBase admin UI: `http://localhost:8080/_/`

---

## Environment Variables

These are baked into the frontend at `docker build` time. Override them by editing the `ENV` lines in the Dockerfile before building.

| Variable | Default | Description |
|---|---|---|
| `VITE_APP_TITLE` | `Process model - Supply Chain Resilience` | Browser `<title>` shown in the tab. |
| `VITE_API_URL` | `/` | PocketBase API base URL. Only needed when the frontend and backend are served from different origins. |

See `.env.example` for a template useful when running the frontend dev server locally.

> **Note:** Never commit `.env` files. They are already in `.gitignore`.

---

## First-Time Setup

### Create admin account

1. Navigate to `http://localhost:8080/_/`
2. Follow the PocketBase setup wizard to create your superuser account.

### Verify seeded data

On first startup, database migrations run automatically. This includes a seed migration that populates the `questions` collection with multilingual vulnerability and capability schemas.

To verify:

1. Log into the admin UI at `http://localhost:8080/_/`
2. Go to **Collections → questions**
3. You should see a record named `v1_multilingual`.

If the record is missing, check container logs:

```bash
docker logs dk-if-open
```

---

## How It Works

### Sessions

An admin creates a **session** via the admin panel. On creation the backend automatically generates:

- 10 participant accounts (`{session-name}-n1` … `{session-name}-n10`) each with a random password and a login token
- 1 session manager account (`{session-name}-s`) with a random password

Participants log in using either their password or their single-use token. Submitting a survey extends the token validity by 3 days.

### Custom API Routes

Beyond standard PocketBase endpoints, the backend exposes:

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/auth-set` | `super` role | Returns whether the session manager has set their own password. |
| `POST` | `/auth-set` | `super` role | Sets the session manager's password. |
| `POST` | `/token?u={userId}` | `super` role | Creates a login token for the given user. |

---

## Email (SMTP)

PocketBase has built-in SMTP support for transactional emails (e.g. password reset). Configure it through the admin UI under **Settings → Mail settings**. The container image includes Postfix if a local relay is preferred.

---

## Development (without Docker)

### Backend

Requires Go ≥ 1.26.3.

```bash
cd back
go run . serve
```

PocketBase admin UI: `http://127.0.0.1:8090/_/`

### Frontend

Requires Node.js ≥ 20.

```bash
cd web
cp ../.env.example .env.local
# Edit .env.local and set VITE_API_URL=http://127.0.0.1:8090
npm install
npm run dev
```

---

## Maintenance

### Backup

PocketBase data lives in the `dk-if-open-data` Docker volume, mounted at `/pb_data`.

**Backup:**
```bash
docker run --rm \
  -v dk-if-open-data:/source:ro \
  -v "$(pwd)/backups":/backup \
  alpine tar -czf /backup/pb_data_$(date +%Y%m%d_%H%M%S).tar.gz -C /source .
```

**Restore:**
```bash
docker stop dk-if-open
docker run --rm \
  -v dk-if-open-data:/target \
  -v "$(pwd)/backups":/backup \
  alpine sh -c "rm -rf /target/* && tar -xzf /backup/pb_data_YYYYMMDD_HHMMSS.tar.gz -C /target"
docker start dk-if-open
```

### Viewing Logs

```bash
docker logs -f dk-if-open
```

---

## Contributing

1. Fork the repository.
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Commit your changes.
4. Open a pull request.
