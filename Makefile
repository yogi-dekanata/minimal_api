run:
	go run cmd/main.go

test:
	go test ./...

build:
	mkdir -p bin
	go build -o bin/app cmd/main.go