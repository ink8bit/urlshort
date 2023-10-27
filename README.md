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
