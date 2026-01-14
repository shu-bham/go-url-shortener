# Go URL Shortener

A simple URL shortener service built with Go.

## Local Setup Instructions

Follow these steps to set up and run the application locally.

### Prerequisites

*   **Go:** [Install Go](https://golang.org/doc/install)
*   **Docker & Docker Compose:** Required for running the MySQL database. [Install Docker](https://docs.docker.com/get-docker/)

### 1. Clone the Repository

```bash
git clone https://github.com/shu-bham/go-url-shortener.git
cd go-url-shortener
```

### 2. Start MySQL Database

Use Docker Compose to spin up the MySQL container:

```bash
docker-compose up -d mysql
```

This will start a MySQL instance accessible at `localhost:3306`.

### 3. Run Database Migrations

First, ensure you have the `migrate` CLI tool installed. If not, install it:
```bash
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Then, run the migration script to set up the database schema.
**Important:** Replace `your_mysql_user`, `your_mysql_password`, and `your_database_name` with the actual credentials from `docker-compose.yml` (or your environment).

```bash
export MYSQL_DSN="your_mysql_user:your_mysql_password@tcp(127.0.0.1:3306)/your_database_name?parseTime=true"
./scripts/migrate.sh up
```

Alternatively, you can manually run the migrations using the `migrate` CLI:

```bash
# Ensure MYSQL_DSN is set as above
migrate -path db/migrations -database "$MYSQL_DSN" up
```

### 4. Configuration

The application can be configured using a `config.yml` file located in the project root, or by setting environment variables.

#### Using `config.yml`

Create a `config.yml` file in the root of the project with the following structure:

```yaml
db:
  dsn: "user:password@tcp(127.0.0.1:3306)/shortener_db?parseTime=true"
logger:
  level: "info"
  format: "json"
server:
  port: "8080"
  domain_name: "http://localhost:8080" # Your application's domain/base URL
```

*Note: Adjust `db.dsn` with the actual credentials and database name from your `docker-compose.yml` or MySQL setup.*



### 5. Run the Application

```bash
go run cmd/shortener/main.go
```

The server should start on the configured port (default 8080).

### 6. Run Tests

```bash
go test ./...
```
