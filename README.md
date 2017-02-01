### Download go
wget https://storage.googleapis.com/golang/go1.7.5.darwin-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.7.5.darwin-amd64.tar.gz

### Download the gRPC package
go get -u google.golang.org/grpc

### Downlaod the protobuf packages
go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go

### Compile the source
protoc -i ./pb ./pb/messages.proto --go_out=plugins=grpc:./src 

