# Backend — Supply Chain Resilience

Go backend built on [PocketBase](https://pocketbase.io/), extended with custom business logic for session and survey management.

## Requirements

- Go ≥ 1.26.3

## Running in development

```bash
go run . serve
```

PocketBase admin UI: `http://127.0.0.1:8090/_/`

## Build

```bash
CGO_ENABLED=0 go build -o backend
./backend serve
```

## Custom behaviour

### Session creation

When a session record is created, the backend automatically generates:

- 10 participant accounts (`{session-name}-n1` … `{session-name}-n10`) with random passwords and login tokens
- 1 session manager account (`{session-name}-s`) with a random password

### Token authentication

Participants can authenticate using a time-limited token (stored in the `tokens` collection) instead of a password. Submitting a survey extends all active tokens for that user by 3 days.

### Custom routes

| Method | Path | Auth required | Description |
|---|---|---|---|
| `GET` | `/auth-set` | `super` role | Check whether the session manager has set their own password. |
| `POST` | `/auth-set` | `super` role | Set the session manager's password. |
| `POST` | `/token?u={userId}` | `super` role | Create a login token for a given user. |

## Migrations

Migration files in `migrations/` run automatically on startup. No manual steps required.
