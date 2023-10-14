# Install Go plugins for the protocol compiler

$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Create the chatserver directory

$ mkdir chatserver

# Generating the go code from the proto file

$ protoc --go-grpc_out=require_unimplemented_servers=false:. --go_out=. chat.proto
