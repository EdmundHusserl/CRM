# CRM

## Prerequisites

1. [Install go](https://go.dev/doc/install)


## Getting started

```sh
$ go mod tidy
$ go test -race ./...
$ go run cmd/main.go -h

Usage of /home/user/.cache/go-build/76/main:
  -db string
        DB provider: in-memory|psql (default "in-memory")
  -port int
        Server port (default 3000)
```


## Deploy using docker compose
```sh
$ docker compose up
```


 ## Environment Variables for Database Connection

The following environment variables are used to configure the database connection.

| Variable     | Default Value | Description |
|-------------|--------------|-------------|
| `DB_HOST`   | `localhost`  | The hostname or IP address of the database server. |
| `DB_PORT`   | `5432`       | The port number on which the database is running. |
| `DB_USER`   | `postgres`   | The username used to authenticate with the database. |
| `DB_PASSWORD` | nil        | The password for the database user (must be set manually). |
|
