## Handling HTTP requests with Automi Streams

This example shows how to create an HTTP server that uses
Automi to process incoming request data and return a response
to the client.

```go
	http.HandleFunc(
		"/",
		func(resp http.ResponseWriter, req *http.Request) {
			resp.Header().Add("Content-Type", "text/html")
			resp.WriteHeader(http.StatusOK)

			// setup new stream with HTTP body as source
			strm := stream.From(sources.Reader(req.Body))
			strm.Run(
				exec.Execute(func(_ context.Context, data []byte) string {
					return base64.StdEncoding.EncodeToString(data)
				}),
			)

			// route result into response
			strm.Into(sinks.Writer[[]byte](resp))

			// run the stream
			if err := <-strm.Open(req.Context()); err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				slog.Error("Stream failed to open", "error", err)
			}
		},
	)
```

The example sets up a request handler with a simple Automi stream
that accepts the incoming data and process it.

#### Testing
Start the server:

```
go run .
```

In a separate terminal, you can use `curl` to send data to the server:

```
curl -d "Hello World!" http://127.0.0.1:4040
```