Automi Streaming and gRPC Streaming
===================================

Because Automi streams use standard built-in Go types, they can be used to stream data to and from
any sources that support native Go type. This example shows how to setup Automi to stream data from 
gRPC streaming servers and clients.

The example in this directory shows how that can be done:

- [server.go](./server.go) - is a gRPC server that emits time on a gRPC client stream.
- [client.go](./client.go) - is a gRPC client uses Automi to stream time values from the gRPC server.
- [protobuf/time.proto](./protobuf/time.proto) - Protobuf file for gPRC service

## Prerequisites
* Download and install the Protocol Buffers compiler `protoc` - https://grpc.io/docs/protoc-installation/
* Install `protoc`'s Go plugins for protobuf and grpc ( see instructions [here](https://grpc.io/docs/languages/go/quickstart/)):
  ```
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```
* Ensure $GOPATH is added to your system's PATH

## Building
First, generate the gRPC Go client/server files with `protoc`.  From within the
`protobuf` directory, run the `protoc` command:

```
$> cd protobuf
$> protoc --go_out=. --go_opt=paths=source_relative  \
   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
   ./time.proto
```
The previous command should generate files:
```
- time_grpc.pb.go  
- time.pb.go
```
## Running the examples
Run the server:
```
go run ./server/server.go
```

In a different terminal, run the client:
```
go run ./client/client.go
```

You should see the time streamed from Automi printed:

```
2025/02/16 23:55:59 INFO Retrieving time stream from server
2025-02-16 23:55:59 +0000 UTC
2025-02-16 23:56:02 +0000 UTC
2025-02-16 23:56:05 +0000 UTC
```
