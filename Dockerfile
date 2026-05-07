FROM golang:1.25-alpine AS builder

RUN apk add --no-cache \
    git \
    ffmpeg \
    gcc \
    musl-dev \
    sqlite-dev

WORKDIR /app

COPY src/go.mod src/go.sum ./
RUN CGO_ENABLED=1 go mod download

COPY src/ ./

ENV CGO_ENABLED=1
ENV GOOS=linux
RUN go build -v -ldflags="-w -s" -o whatsapp .

FROM alpine:latest

RUN apk add --no-cache \
    ffmpeg \
    ca-certificates \
    sqlite-libs \
    tzdata

WORKDIR /app

COPY --from=builder /app/whatsapp .

RUN mkdir -p /app/storages && chmod 777 /app/storages

EXPOSE 10000

CMD ["./whatsapp", "rest", "--port=10000"]
