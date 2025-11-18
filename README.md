# TODO List Backend

Minimal REST API for managing todo items, written in Go.

## Quick start
1) Provision PostgreSQL (anywhere reachable from the app). Create a database (default name: `todos`).
2) Copy `.env.example` to `.env` and update the DB connection values.
3) Run locally:
```bash
go run ./cmd/server
```

### Environment
- `PORT` (default `8080`)
- `DB_HOST` (default `localhost`)
- `DB_PORT` (default `5432`)
- `DB_USER` (default `postgres`)
- `DB_PASSWORD` (required if your DB needs auth)
- `DB_NAME` (default `todos`)
- `DB_SSLMODE` (default `disable`; set to `require` when using TLS)

The server will auto-create a `todos` table if it does not exist.

## API
- `GET /healthz` → `{"status":"ok"}`
- `GET /todos` → list items
- `POST /todos` with `{"title":"Buy milk","completed":false}` → create
- `GET /todos/{id}` → fetch single item
- `PUT /todos/{id}` with any of `title`, `completed` → update
- `DELETE /todos/{id}` → remove

Example create then toggle complete:
```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Try the new API"}'

curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"completed":true}'
```

## Docker
Build and run:
```bash
docker build -t todo-list-be .
docker run --rm -p 8080:8080 --env-file .env todo-list-be
```

## Kubernetes
Helm chart under `deploy/todo-list-be` deploys the API (listens on port 8080). Configure database access through the `env.*` values (DB host/port/user/password/name/sslmode) and override the image repository/tag as needed.

## Make targets
- `make build` – build dev Docker image (`$(DOCKERHUB_USER)/todo-list-be:$(DEV_VERSION)`)
- `make test` – run Go tests
- `make push` – push the dev image
- `make build-prod` / `make push-prod` – production image & push
