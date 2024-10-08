# Dockerfile para el microservicio account
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o account_service .

# Imagen final
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/account_service .

# Definir variables de entorno
ENV SPRING_RABBITMQ_URL=amqp://admin:admin@rabbitmq:5672/
ENV SPRING_RABBITMQ_HOST=rabbitmq
ENV SPRING_RABBITMQ_PORT=5672
ENV SPRING_RABBITMQ_USERNAME=admin
ENV SPRING_RABBITMQ_PASSWORD=admin

EXPOSE 8082

CMD ["./account_service"]
