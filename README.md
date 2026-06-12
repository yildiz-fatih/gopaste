# gopaste

GoPaste is a minimal, self-hostable pastebin built with Go, PostgreSQL and Redis.

## Quick Start

```bash
cp .env.example .env
# open .env and change the password values

make up # runs on http://localhost:8080
```

## Development

```bash
cp .env.example .env
# open .env and change the password values
 
make dev

make run # runs on http://localhost:8080
```

## Testing

```bash
make test
```
