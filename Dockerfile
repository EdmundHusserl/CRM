FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Run unit tests with race condition checks
RUN go test -race -cover ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C cmd/ -o /go/bin/app

FROM alpine:latest
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
COPY --from=builder /go/bin/app /app

ENTRYPOINT ["/app"]
