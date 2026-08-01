# ---- etapa de build ----
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd

# ---- etapa final ----
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

RUN adduser -D -H appuser
USER appuser

COPY --from=builder /server /server

EXPOSE 8080

ENTRYPOINT ["/server"]