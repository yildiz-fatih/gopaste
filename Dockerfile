# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /gopaste

# Run stage
FROM alpine:latest

COPY --from=builder /gopaste /gopaste
COPY --from=builder /app/views /views
COPY --from=builder /app/static /static
COPY --from=builder /app/help.md /help.md

CMD ["/gopaste"]
