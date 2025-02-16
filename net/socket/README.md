## Processing socket data with Automi

This example shows how to create a socket server that uses
Automi to process incoming bytes from the client, process the data, 
and return a response.

This is a simple echo server that illustrate how this would work.

```go
func main() {
	addr := ":4040"
	ln, err := net.Listen("tcp", addr)
	defer ln.Close()

	conn, err := ln.Accept()
	defer conn.Close()

	// Automi is used to stream from the connection an io.Reader source.
	// The stream is transformed then the result is routed to an io.Writer sink.
	strm := stream.From(sources.Reader(conn))
	strm.Run(
		exec.Map(func(_ context.Context, chunk []byte) string {
			return strings.ToUpper(string(chunk))
		}),
	)
	strm.Into(sinks.Writer[string](conn))

	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}
```

You can tests by first starting the socket server:

```
go run .
```

In a separate terminal, use `netcat` (or similar) to connect:

```
nc localhost 4040
Hello World
HELLO WORLD
```