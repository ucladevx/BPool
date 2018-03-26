# First stage get deps and build
FROM golang:1.9 AS builder

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com/ucladevx/BPool
WORKDIR /go/src/github.com/ucladevx/BPool

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/BPool/main.go

## Second stage run the app
FROM alpine:3.7

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /root

COPY --from=builder /go/src/github.com/ucladevx/BPool/BPool.out .

CMD ["./BPool"]