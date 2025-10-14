# Fuzzy Builder

## Backend

- go build -o bin/api cmd/api/main.go
- DATABASE_URL=... JWT_SECRET=... ./bin/api

## Migrations

- go build -o bin/migrator cmd/migrator/main.go
- DATABASE_URL=... ./bin/migrator

## Docker Compose

- docker compose build
- docker compose up -d
- API: http://localhost:8080
- Frontend: http://localhost:5173

## Frontend

See `frontend/README.md`. Configure `.env`:

```
VITE_API_BASE=http://localhost:8080
```
