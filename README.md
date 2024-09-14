# Snippetbox

Web applicaton that enables writting code snippets or other forms of texts and stores them in a PostgreSQL database together with their title and a TTL.

## Quick start:

1. Start PostgreSQL with docker compose: `docker-compose up`
2. Seed the database withe the models found under the `schemas` directory using your DBA of choice
3. Start the Go server: `go run ./cmd/web`

## Technologies
Web application fully written in Go v1.22

## Notes
Web application based on the application described in the [Let's go](https://lets-go.alexedwards.net/) book by Alex Edwards.
