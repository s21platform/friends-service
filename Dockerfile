FROM golang:1.22-alpine as builder
RUN apk update && apk add --no-cache make

WORKDIR /usr/src/service
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -o build/main cmd/service/main.go
#RUN go build -o build/kafka cmd/workers/notification/main.go TODO:понять почему тут ошибка

FROM alpine:latest

WORKDIR /app

COPY --from=builder /usr/src/service/build/main .
#COPY --from=builder /usr/src/service/build/kafka . TODO: поправить тут тоже

CMD ["/app/main","/app/kafka"]