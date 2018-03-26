default: build

.PHONY: config gen-config

run:
	go run ./cmd/BPool/main.go

build:
	go build -o BPool.out ./cmd/BPool/main.go 

gen-config:
	gpg -c ./config/config.dev.json

config:
	gpg ./config/config.dev.json.gpg