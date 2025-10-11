run:
    go run cmd/app/main.go
fmt:
    golangci-lint fmt
docs:
    swag init -g cmd/app/main.go
tests:
    cd tests && uv run main.py