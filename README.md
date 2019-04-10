# github.com/e2fyi/minio-web
[![godoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square "godoc")](https://godoc.org/github.com/e2fyi/minio-web/pkg)

A web server proxy for any S3-compatible storage.


```bash
# starts a minio server
./scripts/start-minio-server.sh

# starts minio-web service
go run cmd/minio-web.go
```

Alternatively, with docker compose:
```bash
# you can access the minio-store at localhost:9000
# and the minio-web at localhost:8080
docker-compose up -d
```

docker image:
```bash
# build locally
docker build -t e2fyi/minio-web:latest .

# pull from dockerhub
docker pull e2fyi/minio-web:latest

# run docker container
docker run --rm -ti \
    -p 8080:8080 \
    -e MINIO_ENDPOINT=s3.amazonaws.com \
    -e MINIO_PORT=443 \
    -e MINIO_SSL=true \
    -e MINIO_ACCESS_KEY_ID=ABCD \
    -e MINIO_SECRET_ACCESS_KEY=EFGH \
    -e MINIO_BUCKET=FooBar \
    e2fyi/minio-web:latest
```


## GoDoc

- [minio-web](https://godoc.org/github.com/e2fyi/minio-web/cmd)
- [github.com/e2fyi/minio-web/pkg](https://godoc.org/github.com/e2fyi/minio-web/pkg)