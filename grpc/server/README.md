## Time server
The server code is simple. It exposes a method `GetTimeStream` which returns a stream of time stamps:

```go
func (t *timeServer) GetTimeStream(req *pb.TimeRequest, stream pb.TimeService_GetTimeStreamServer) error {
	delay := time.Second
	interval := req.GetInterval()
	if interval > 0 {
		delay = time.Millisecond * time.Duration(interval)
	}

	slog.Info("Creating time stream", "delay", delay, "interval", interval)

	buf := make([]byte, 8)
	for {
		binary.BigEndian.PutUint64(buf, uint64(time.Now().Unix()))
		stream.Send(&pb.Time{Value: buf})
		time.Sleep(delay)
	}
}
```