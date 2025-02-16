package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
)

func main() {
	addr := ":4040"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("Creating listener failed", "error", err)
		os.Exit(1)
	}
	defer ln.Close()

	slog.Info("Service started", "port", addr)

	conn, err := ln.Accept()
	if err != nil {
		slog.Error("Unable to connect to client", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	slog.Info("Connected client", "address", conn.RemoteAddr())

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
