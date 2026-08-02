FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY pkg/ ./pkg/
COPY services/reading/ ./services/reading/

WORKDIR /app/services/reading
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/server ./cmd/main.go

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

COPY --from=builder /app/bin/server /server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/server"]
