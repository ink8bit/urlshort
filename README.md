# urlshort

## Getting started

1. Install packages:

    ```console
    go get
    ```

1. Run application:

    ```console
    go run cmd/shortener/main.go
    ```

1. Open your browser:

    ```console
    open http://localhost:8080
    ```

## Testing

Run all tests:

```console
go test -v ./...
```

Check test coverage:

```console
go test -coverprofile c.out ./...
```

View test coverage in browser:

```console
go tool cover -html=c.out
```

## PostgreSQL

> This installation steps provided for macOS only.

### Installation

> View [docs](https://wiki.postgresql.org/wiki/Homebrew)

```console
brew install postgres
```

### Start and stop postgres

```console
brew services start postgresql
brew services stop postgresql
```

### Setup

```console
psql postgres
```

#### View users

```console
\du
```

#### Create a new user

```console
CREATE ROLE postgres WITH LOGIN PASSWORD '<YOUR_PWD>';
ALTER ROLE postgres CREATEDB;
```

#### Create a database

```console
CREATE DATABASE shorturls;
```
