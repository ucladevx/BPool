default: build

.PHONY: config gen-config test

run:
	go run ./cmd/BPool/main.go

build:
	go build -o BPool.out ./cmd/BPool/main.go 

test:
	go test -v ./...

gen-config:
	gpg -c ./config/config.dev.json

config:
	gpg ./config/config.dev.json.gpg