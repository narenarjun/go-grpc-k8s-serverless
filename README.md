# <div align="center"> Golang GRPC blog server on deployed kubernetes made serverless with Knative</div>



The kubernetes deployment files are in the [`./kubernetes/server-deployment`]("./kubernetes/server-deployment") folder.

To create a client for this server, we can use this proto file : [`./blog-server/blogpb/blog.proto`]("./blog-server/blogpb/blog.proto").

The grpc server is deployed on the [`civo`](civo.com) k3s kubernetes platform. Here is the server url: [grpc-blog-server.grpc-blog.e20b4706-9ba3-4496-a857-b8b531dd5a38.k8s.civo.com]("grpc-blog-server.grpc-blog.e20b4706-9ba3-4496-a857-b8b531dd5a38.k8s.civo.com")

## ✨✨ Container Image:

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