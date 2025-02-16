package main

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	autoctx "github.com/vladimirvivien/automi/api/context"
	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"
	"github.com/vladimirvivien/automi/stream"
	"github.com/vladimirvivien/automi/log"
)

type walkInfo struct {
	path string
	err  error
}

// Demostrates stream runtime logging. Uses the examples/md5 as basis.
// To run: go run md5.go -p ./..
func emitPathsFor(root string) <-chan walkInfo {
	paths := make(chan walkInfo)
	go func() {
		defer close(paths)
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if !info.Mode().IsRegular() {
				return nil
			}
			paths <- walkInfo{path, err}
			return nil
		})
	}()
	return paths
}

func main() {
	var rootPath string
	flag.StringVar(&rootPath, "p", "./", "Root path to start scanning")
	flag.Parse()

	ctx := context.Background()
	stream := stream.From(sources.Chan(emitPathsFor(rootPath)))
	stream.WithLogSink(sinks.SlogJSON(slog.LevelDebug))
	stream.Log(ctx, log.LogInfo("Walking path", slog.String("path", rootPath)))

	stream.Run(
		// filter out errored walk results
		exec.Filter(func(ctx context.Context, info walkInfo) bool {
			return info.err == nil
		}),

		// map tuple walkInfo -> string
		exec.Map(func(ctx context.Context, info walkInfo) string {
			autoctx.LogF(ctx, log.LogInfo("selecting path", slog.String("path", info.path)))
			return info.path
		}),

		// mapping file content to md5 sum
		exec.Map(func(ctx context.Context, filePath string) [3]any {
			autoctx.LogF(ctx, log.LogInfo("generating md5 sum", slog.String("path", filePath)))
			data, err := os.ReadFile(filePath)
			sum := md5.Sum(data)
			return [3]any{filePath, sum, err}
		}),
	)

	// sink the result
	stream.Into(sinks.Func(func(items [3]any) error {
		file := items[0].(string)
		md5Sum := items[1].([md5.Size]byte)
		fmt.Printf("file %-64s md5 (%-16x)\n", file, md5Sum)
		return nil
	}))

	// open the stream
	if err := <-stream.Open(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
