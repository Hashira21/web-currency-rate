FROM golang:1.23.1-alpine as builder

RUN apk add --no-cache git

WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/currency-rate/main.go

FROM alpine:3.18
RUN apk update && apk add tzdata
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

EXPOSE 8080

COPY --from=builder /src/app .
COPY configs configs
COPY backup.sql backup.sql
COPY wait-for-postgres.sh wait-for-postgres.sh

# install psql
RUN apk add postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

CMD ["/app"]

# Фронтенд (новый Dockerfile в директории frontend/Dockerfile)
# frontend/Dockerfile
FROM nginx:alpine

# Копируем статические файлы из папки frontend
COPY frontend/index.html /usr/share/nginx/html/
COPY frontend/css/ /usr/share/nginx/html/css/
COPY frontend/script.js /usr/share/nginx/html/

# Копируем конфиг Nginx из frontend/nginx
COPY frontend/nginx/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80