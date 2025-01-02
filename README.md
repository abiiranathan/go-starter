# go-starter

An opinionated starter template for building robust web applications with Go.

### Tech Stack:

- [Rex Router](https://github.com/abiiranathan/rex) A fast and flexible HTTP router based on golang's Go1.22 enhanced routing.
- [sqlc](https://github.com/sqlc-dev/sqlc) for managing the the postgres database (with pgx/v5 driver) and [pgxpool](https://github.com/jackc/pgx/v5/pgxpool) for connection pooling.
- [go-redis](github.com/redis/go-redis/v9) for redis caching.
- [golang migrate](https://github.com/golang-migrate/migrate) for database migrations.
- [templ](https://github.com/a-h/templ) for template management
- [tailwindcss](https://tailwindcss.com/) for styling.
- [air](https://github.com/air-verse/air) for live reload.
- [.env or config.yaml](./config.yml) for the configuration

## Getting Started

### Prerequisites

Install the required tools:

```bash
make install
```

This will install sqlc and golang-migrate.

Then install redis and postgresql.

```bash
sudo apt install redis-server postgresql postgresql-contrib
```

### Setup

Create a `.env` file in the root directory and add the following environment variables:

```bash
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
PORT=8080
SECRET_KEY=secret
DEBUG=true
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5174
```

### Build

```bash
make build
```

### Workflow

```bash
make dev
```

This will start the server with live reload.

In another terminal, run the following command to start tailwindcss watcher:

```bash
make watch
```

In another terminal, run the following command to run the `templ` compiler:

```bash
templ generate --watch
```

### Migrations

To create a new migration, run the following command:

```bash
make create
```

To apply the migrations, run the following command:

```bash
make up
```

### Generate SQLC

To generate the sqlc code, run the following command:

```bash
make generate
```

See the [Makefile](./Makefile) for more commands.
