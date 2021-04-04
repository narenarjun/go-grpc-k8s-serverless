# Golang GRPC blog server on kubernetes made serverless with Knative 

<!-- to generate a grpc code, `protoc-gen-go-grpc` plugin must be used instead of `protoc-gen-go` plugin


so the generate command is 

```
protoc --go-grpc_out=./ --go_out=./ ./blogpb/blog.proto
```

to generate the grpc part `--go-grpc_out` option is added and to generate the go part `--go_out` option is added

inside the blog.proto file, the option go_package value must follow this convention 
https://developers.google.com/protocol-buffers/docs/reference/go-generated#package
 -->