### Download go
wget https://storage.googleapis.com/golang/go1.7.5.darwin-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.7.5.darwin-amd64.tar.gz

### Download Protocol Buffer Compiler
. DL correct file from https://github.com/google/protobuf/releases
. Update path to point to 'protoc' binary

### Download the gRPC package
go get -u google.golang.org/grpc

### Downlaod the protobuf packages
go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go

### Compile the source
protoc -I ./pb ./pb/messages.proto --go_out=plugins=grpc:./pb 

### Use gRPC_CLI and Server Reflection to debug
had to use https://github.com/grpc/homebrew-grpc to build (and add a few libs)
then ran make grpc_cli
https://chromium.googlesource.com/external/github.com/grpc/grpc-go/+/refs/heads/master/Documentation/server-reflection-tutorial.md

### go generate?
https://blog.golang.org/generate

### gRPC tracing
https://rakyll.org/grpc-trace/

