FROM golang:1.18 AS builder

COPY main.go go.mod go.sum /build/
WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch
COPY --from=builder /build/librivox-sorter /

# Grab some current ca-certs
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT  ["/librivox-sorter"]
