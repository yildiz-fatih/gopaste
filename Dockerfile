# Build stage
FROM golang:1.26.5-alpine3.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /gopaste ./cmd/web

# Run stage
FROM alpine:3.24.1

COPY --from=builder /gopaste /gopaste
COPY --from=builder /app/cmd/web/views /views
COPY --from=builder /app/static /static
COPY --from=builder /app/help.md /help.md

CMD ["/gopaste"]
