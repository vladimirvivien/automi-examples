## Aggregate - Group by Struct Field

This  example shows how to aggregate streamed data items stored in Go struct types.
First the code declares a struct type for the data and sets up a slice of the struct
as the data source for the stream:

```go
type log struct{ Event, Src, Device, Result string }
source := sources.Slice([]log{ ... }
```

Next, the code sets up the stream to, first, batch each streamed item
into a slice of logs `[]log`. Then, the batched items are grouped using
the value of struct field `Src` and the result is collected into a map
of type `mam[any][]log`.

```go
	// Create new stream from source
	strm := stream.From(source)

	// Setup stream operations
	strm.Run(
		// Batch each incoming `log` items into --> `[]log`
		window.Batch[log](),

		// Group incoming `[]log` by struct field `Src` into --> map[any][]log
		exec.GroupByStructField[[]log]("Src"),
	)

	// Define sink opertor to handle items
	strm.Into(sinks.Func(func(items map[any][]log) error {
		for key, item := range items {
			fmt.Println("Group:", key)
			fmt.Println("-----")
			for _, log := range item {
				fmt.Printf("log = %v\n", log)
			}
		}
		return nil
	}))
```