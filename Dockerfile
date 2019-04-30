#################################
FROM golang:1.12-alpine as builder

WORKDIR /goapp

RUN apk --no-cache add git

COPY . .

RUN go build .

#################################
FROM alpine:latest as dist

WORKDIR /goapp

RUN apk --no-cache add su-exec ca-certificates

COPY --from=builder /goapp/minio-web  /goapp/minio-web
COPY assets/ assets/
COPY configs/config.json config.json

ENV CONFIG_FILE_PATH="/goapp/config.json"
EXPOSE 8080 

ENTRYPOINT [ "su-exec", "goapp:1000" ]

ENV EXT_DEFAULTHTML=index.html
ENV EXT_FAVICON=assets/favicon.ico
ENV EXT_MARKDOWNTEMPLATE=assets/md-template.html

CMD ["/goapp/minio-web"]