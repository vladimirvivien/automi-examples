Automi Logging
===============

This example (builds on the md5 example) shows how to setup Automi to emitt log events during stream operations.
It configures an slog.JSNONHandler as log sink to display log events on stdout.

```go
stream.WithLogSink(sinks.SlogJSON(slog.LevelDebug))
```

The example also uses the Automi context to extract the logger to log events
during the execution of an operation:

```go
exec.Map(func(ctx context.Context, info walkInfo) string {
	autoctx.LogF(ctx, util.LogInfo("selecting path", slog.String("path", info.path)))
	return info.path
})
```
 