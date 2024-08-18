FROM golang:1.22 as builder

WORKDIR /usr/src/service
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

#RUN go build -o build/kafka cmd/workers/notification/main.go todo: fix
RUN go build -o build/main cmd/service/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /usr/src/service/build/main .
COPY --from=builder /usr/src/service/scripts ./scripts
#COPY --from=builder /usr/src/service/build/kafka . todo:fix

#CMD ["/app/main","/app/kafka"]


CMD ["./app/main"]