# Golang GRPC blog server on kubernetes made serverless with Knative 



## Container Image

The app is containerized and hosted in the [Quay Container registry](https://quay.io/repository/narendev/blogserver?tab=tags).   

The app is built by multistage build to reduce the final image size.     

To get the image:

with docker:

```bash

docker pull quay.io/narendev/blogserver

```


with podman:

```bash

podman pull quay.io/narendev/blogserver

```


<!-- to generate a grpc code, `protoc-gen-go-grpc` plugin must be used instead of `protoc-gen-go` plugin


so the generate command is 

```
protoc --go-grpc_out=./ --go_out=./ ./blogpb/blog.proto
```

to generate the grpc part `--go-grpc_out` option is added and to generate the go part `--go_out` option is added

inside the blog.proto file, the option go_package value must follow this convention 
https://developers.google.com/protocol-buffers/docs/reference/go-generated#package
 -->