# Сначала создаем бинарный файл на основе образа golang:latest для amd64
FROM --platform=linux/amd64 golang:latest AS builder-amd64

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOARCH=amd64 go build -o main .

# Теперь создаем бинарный файл для arm64
FROM --platform=linux/arm64 golang:latest AS builder-arm64

WORKDIR /app

COPY . .
COPY ok_server.crt minka/goproj/Back-End-API/ok_server.crt
COPY ok_server.key ok_server.key

RUN CGO_ENABLED=0 GOARCH=arm64 go build -o main .

# Теперь создаем конечный образ с использованием образа ubuntu
FROM ubuntu:latest

WORKDIR /app

# Копируем бинарный файл в зависимости от архитектуры
COPY --from=builder-amd64 /app/main /app/main-amd64
COPY --from=builder-arm64 /app/main /app/main-arm64

# Устанавливаем необходимые зависимости
RUN apt-get update && apt-get install -y \
    libc6-dev \
    && rm -rf /var/lib/apt/lists/*

# В зависимости от платформы запускаем нужный бинарный файл
CMD ["/bin/sh", "-c", "if [ \"$(uname -m)\" = \"x86_64\" ]; then /app/main-amd64; else /app/main-arm64; fi"]
