# gopaste

GoPaste is a minimal, self-hostable pastebin built with Go, PostgreSQL and Redis.

## Quick Start

```bash
cp .env.example .env
# open .env and change the password values

docker compose up # runs on http://localhost:8080
```

## Development

```bash
cp .env.example .env
# open .env and change the password values
 
docker compose up db redis

export DATABASE_URL=<DATABASE_URL from your .env>
migrate -path migrations -database $DATABASE_URL up

go run . # runs on http://localhost:8080
```
