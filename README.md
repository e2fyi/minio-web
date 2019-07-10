# github.com/e2fyi/minio-web
[![godoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square "godoc")](https://godoc.org/github.com/e2fyi/minio-web/pkg) [![dockerhub](https://img.shields.io/badge/dockerhub-e2fyi%2Fminio--web-5272B4.svg?style=flat-square "dockerhub")](https://hub.docker.com/r/e2fyi/minio-web)

A web server proxy for any S3-compatible storage.

## Quickstart

### Environment variables

```bash
# where config file is located
CONFIG_FILE_PATH=configs/config.json

# port for minio-web to listen to
SERVER_PORT=8080
# path to ssl cert and key files
SERVER_SSL_CERT=
SERVER_SLL_KEY=

# endpoint to call for the s3 compatible storage
MINIO_ENDPOINT=s3.amazonaws.com
# access key and secret key
MINIO_ACCESSKEY=
MINIO_SECRETKEY=
# ssl when calling endpoint
MINIO_SECURE=true
# aws s3 bucket region (optional)
MINIO_REGION=

# Extensions #
# bucket to serve if provided (http://minio-web/abc => endpoint/bucketname/abc)
# if not provided (http://minio-web/abc/efg => endpoint/abc/efg) where abc is the bucket
EXT_BUCKETNAME=

# if provided a default index file is return 
# i.e http://minio-web/abc/ => http://minio-web/abc/index.html
EXT_DEFAULTHTML=index.html

# if provided, returns a default favicon if backend does not have one.
EXT_FAVICON=assets/favicon.ico

# if set, list the folders inside a folder
EXT_LISTFOLDER=true
# if set, list all objects inside a folder, otherwise only list folders
EXT_LISTFOLDEROBJECTS=false

# if provided, renders any markdown resources as HTML with the template.
# template MUST have a placeholder {{ .Content }}
EXT_MARKDOWNTEMPLATE=assets/md-template.html
```

### Config file

```json
{
    "server": {
        "port": 8080,
        "ssl": {
            "cert": "",
            "key": ""
        }
    },
    "minio": {
        "endpoint": "s3.amazonaws.com",
        "accesskey": "",
        "secretkey": "",
        "secure": false ,
        "region": ""       
    },
    "ext": {
        "bucketname": "",
        "defaulthtml": "index.html",
        "favicon": "assets/favicon.ico",
        "markdowntemplate": "assets/md-template.html",
        "listfolder": true,
        "listfolderobjects": false
    }
}
```

### Run demo locally
```bash
# starts a minio server
./scripts/start-minio-server.sh

# starts minio-web service
go run .
```

Alternatively, with docker compose:
```bash
# you can access the minio-store at localhost:9000
# and the minio-web at localhost:8080
docker-compose up -d
```

### Docker image
```bash
# build locally
docker build -t e2fyi/minio-web:latest .

# pull from dockerhub
docker pull e2fyi/minio-web:latest

# run docker container
docker run --rm -ti \
    -p 8080:8080 \
    --env-file .envfile \
    e2fyi/minio-web:latest
```

## Kubernetes deployment (Kustomize)
[kustomize](https://github.com/kubernetes-sigs/kustomize) k8s manifest for 
`minio-web` can be found in [manifest/](./manifest).

## GoDoc

- [minio-web](https://godoc.org/github.com/e2fyi/minio-web/)
- [github.com/e2fyi/minio-web/pkg/app](https://godoc.org/github.com/e2fyi/minio-web/pkg/app)
- [github.com/e2fyi/minio-web/pkg/core](https://godoc.org/github.com/e2fyi/minio-web/pkg/core)
- [github.com/e2fyi/minio-web/pkg/minio](https://godoc.org/github.com/e2fyi/minio-web/pkg/minio)
- [github.com/e2fyi/minio-web/pkg/ext](https://godoc.org/github.com/e2fyi/minio-web/pkg/ext)