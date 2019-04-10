#!/bin/bash
docker run --rm -ti \
    -v $PWD/tmp/minio:/data \
    -p 9000:9000 \
    -e "MINIO_ACCESS_KEY=minio" \
    -e "MINIO_SECRET_KEY=minio123" \
    minio/minio server /data