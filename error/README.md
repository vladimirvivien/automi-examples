Automi Error Handling
=====================

This example shows how to setup Automi to generate error events from stream operations
and use type `api.StreamResult` to generate operation results that can contain data
or error to be reported.

The example also sets a default sink logger to capture stream events:

```go
stream.WithLogSink(sinks.SlogJSON(slog.LevelDebug))
```

 