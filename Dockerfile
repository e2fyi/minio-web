#################################
FROM golang:1.12-alpine as builder

WORKDIR /goapp

RUN apk --no-cache add git

COPY . .

RUN go build cmd/minio-web.go

#################################
FROM alpine:latest as dist

WORKDIR /goapp

RUN apk --no-cache add su-exec

COPY --from=builder /goapp/minio-web  /goapp/minio-web
COPY assets/ assets/
COPY configs/config.json config.json

ENV CONFIG_FILE_PATH="/goapp/config.json"
EXPOSE 8080 

ENTRYPOINT [ "su-exec", "goapp:1000" ]

CMD ["/goapp/minio-web"]