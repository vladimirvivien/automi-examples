## The gRPC client

The gRPC client does the followings:
* Connects to the server and retrieves its stream
* Sets up an Automi stream that processes the server's stream data

#### Connect to server

The first thing done in the client is to connec to the service:

```go

conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
if err != nil {
	log.Fatal(err)
}
```
#### Setup data stream
In function `emitStream` the code requests the server's stream as shown below:

```go
timeStream, err := client.GetTimeStream(context.Background(), &pb.TimeRequest{Interval: 3000})
if err != nil {
	log.Fatal(err)
}
```

Next, data from the server's stream is sent to a channel that will be used as an Automi source.

```go
func emitStream(client pb.TimeServiceClient) <-chan []byte {
    source := make(chan []byte)
...
	go func(stream pb.TimeService_GetTimeStreamClient, srcCh chan []byte) {
		for {
			t, err := stream.Recv()
			if err != nil {
				slog.Error("gRPC stream error", "error", err)
				continue
			}
			srcCh <- t.Value
		}
    ...
	}(timeStream, source)
...
    return source
}
```

In the `main` function, the Automi stream is defined to stream data from the channel setup above:

```go
func main() {
	// Setup Automi stream to process time data
	strm := stream.From(sources.Chan(emitStream(client)))
	strm.Run(
		exec.Map(func(ctx context.Context, item []byte) time.Time {
			secs := int64(binary.BigEndian.Uint64(item))
			return time.Unix(int64(secs), 0)
		}),
	)
	strm.Into(
		sinks.Func(func(item interface{}) error {
			time := item.(time.Time)
			fmt.Println(time)
			return nil
		}),
	)
}
```