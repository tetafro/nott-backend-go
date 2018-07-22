# Nott

Markdown notes service with code syntax highlighting.

This repository provides backend written in go.

## Build and run

Define database creds
```sh
export PGUSER=postgres
export PGPASSWORD=postgres
export PGDATABASE=nott
```

Setup PostgreSQL
```sh
docker run -d \
    --name postgres \
    --publish 127.0.0.1:5432:5432 \
    --env "POSTGRES_USER=${PGUSER}" \
    --env "POSTGRES_PASSWORD=${PGPASSWORD}" \
    --env "POSTGRES_DB=${PGDATABASE}" \
    postgres:10
```

Setup database scheme using [migrator](https://github.com/golang-migrate/migrate)
```sh
export DB="postgres://${PGUSER}:${PGPASSWORD}@localhost:5432/${PGDATABASE}?connect_timeout=1&sslmode=disable"
migrate -path migrations -verbose -database ${DB} up
```

Populate with test data
```sh
psql -h localhost ${PGDATABASE} < ./migrations/populate.sql
```

Create and populate config
```sh
cp .env.example .env
source .env
```

Build and run the application
```sh
make build run
```
