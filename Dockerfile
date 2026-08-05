# ---- etapa de build ----
FROM golang:1.25.11-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /migrate ./migrate

# ---- etapa final ----
# ---- etapa final ----
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

RUN adduser -D -H appuser

# Crear el directorio con el dueño correcto ANTES del volumen y del USER
RUN mkdir -p /app/media/products && chown -R appuser:appuser /app/media

USER appuser

COPY --from=builder /server /server
COPY --from=builder /migrate /migrate

EXPOSE 4005

CMD ["/server"]




