# Etapa de compilación
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Etapa de ejecución
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/server ./server

EXPOSE 8080

CMD ["./server"]
