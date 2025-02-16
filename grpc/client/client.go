package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/vladimirvivien/automi/operators/exec"
	"github.com/vladimirvivien/automi/sinks"
	"github.com/vladimirvivien/automi/sources"

	"github.com/vladimirvivien/automi/stream"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/vladimirvivien/automi-examples/grpc/protobuf"
)

const (
	server     = "127.0.0.1"
	serverPort = "50051"
)

// emitStream returns a channel that emits time event from server stream
func emitStream(client pb.TimeServiceClient) <-chan []byte {
	source := make(chan []byte)

	// create gRPC stream source
	slog.Info("Retrieving time stream from server")
	timeStream, err := client.GetTimeStream(context.Background(), &pb.TimeRequest{Interval: 3000})
	if err != nil {
		log.Fatal(err)
	}

	// stram item from gPRC and forward to emitter channel
	go func(stream pb.TimeService_GetTimeStreamClient, srcCh chan []byte) {
		defer close(srcCh)
		for {
			t, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					slog.Info("Server stream closing")
					return // done
				}
				slog.Error("gRPC stream error", "error", err)
				continue
			}
			srcCh <- t.Value
		}

	}(timeStream, source)

	return source
}

func main() {
	serverAddr := net.JoinHostPort(server, serverPort)

	// setup insecure grpc connection
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewTimeServiceClient(conn)

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

	// open the stream
	if err := <-strm.Open(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
