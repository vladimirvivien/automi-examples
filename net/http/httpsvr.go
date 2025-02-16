package main

import (
	"context"
	"encoding/base64"
	"log/slog"
	"net/http"

	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

// HTTP example that returns a base64 encoding of http Body.
// Start server: go run httpsvr.go
// Test with: curl -d "Hello World"  http://127.0.0.1:4040/
func main() {

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

	slog.Info("HTTP server listening on :4040")
	http.ListenAndServe(":4040", nil)
}
