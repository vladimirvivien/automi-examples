## Windowing - Group By Map Key

This example demonstrate how to use Automi to apply `Group By` aggregate operations.
The code sets up a slice of maps as the source for the stream and applies several
stream operations. 

AFter filtering the data, the code batches all stremed items and then groups the 
items using values from field `Devices` in the map.

```go
src := sources.Slice([]map[string]string{ ... }

	// Define new stream from slice source
	strm := stream.From(src)

	// attach a logger sink
	strm.WithLogSink(sinks.SlogText(slog.LevelDebug))

	// define stream operations
	strm.Run(
		// filter out "response" event
		exec.Filter(func(ctx context.Context, e map[string]string) bool {
			return (e["Event"] == "response")
		}),

		// Batch map items into -> []map[string]string
		window.Batch[map[string]string](),

		// Group batched items []map[string]string -> map[string][]map[string]string
		exec.GroupByMapKey[[]map[string]string]("Device"),
	)

	strm.Into(sinks.Func(func(items map[string][]map[string]string) error {
		for batch, items := range items {
			fmt.Println("Batch:", batch)
			fmt.Println("-----")
			for _, item := range items {
				fmt.Printf("%#v\n", item)
			}
		}
		return nil
	}))

```