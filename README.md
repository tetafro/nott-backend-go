# Nott

[![CircleCI](https://circleci.com/gh/tetafro/nott-backend-go.svg?style=shield)](https://circleci.com/gh/tetafro/nott-backend-go)
[![Codecov](https://codecov.io/gh/tetafro/nott-backend-go/branch/master/graph/badge.svg)](https://codecov.io/gh/tetafro/nott-backend-go)
[![Go Report](https://goreportcard.com/badge/github.com/tetafro/nott-backend-go)](https://goreportcard.com/report/github.com/tetafro/nott-backend-go)

Markdown notes service with code syntax highlighting.

This repository provides backend written in go.

## Build and run

Run PostgreSQL
```sh
docker run -d \
    --name postgres-nott \
    --publish 127.0.0.1:5432:5432 \
    --env "POSTGRES_USER=postgres" \
    --env "POSTGRES_PASSWORD=postgres" \
    --env "POSTGRES_DB=nott" \
    postgres:10
```

Create and populate config
```sh
cp .env.example .env
```

Build and run the application
```sh
make build run
```
