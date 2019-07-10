# minio-web/manifest

This is a [kustomize](https://github.com/kubernetes-sigs/kustomize) k8s manifest
for `minio-web`.

`/base` contains the minimal but functional yamls for k8s deployments. You can
either use `kustomize` to modify the yamls or edit them directly.

`/overlay` contains an example overlay with:
- `ambassador` mapping 
- parameterization with `kustomize` (see [`params.env`](./overlay/params.env))