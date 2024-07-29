FROM golang:1.22 as builder
RUN apk update && apk add --no-cache make

WORKDIR /usr/src/service
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -o build/main cmd/service/main.go

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /usr/src/service/build/main .
CMD ["/app/main"]