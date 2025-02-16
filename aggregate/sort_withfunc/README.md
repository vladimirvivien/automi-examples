## Aggregate - Sort with Func

This example shows how to use the `exec.SortWithFunc` operator to sort aggreated items
using a user-specified sort function. First, the code sets up a type, `log`, for the
streamed items and uses a slice of logs `[]log` as the data source.

```go
type log struct{ Event, Src, Device, Result string }
source := sources.Slice([]log{ ... }
```

Next, the code defines a new Automi stream using the source above. It then applies
a filter to select only a speset of the data. The stream has a `Batch` operator that
is used to collect each `log` item into a slice of logs `[]log`. The `exec.SortWithFunc`
function is then used to sort the batched items.

```go
	// Create new stream using source
	strm := stream.From(source)

	// Define stream operators
	strm.Run(
		// Filter some data
		exec.Filter(func(_ context.Context, e log) bool {
			return (e.Event == "response")
		}),

		// Batch the data into --> []log
		window.Batch[log](),

		// Define a sort function for the streamed items
		exec.SortWithFunc[[]log](func(i, j log) int {
			return cmp.Compare(i.Src, j.Src)
		}),
	)

	// Send stream items into a sink handler
	strm.Into(sinks.Func(func(sorted []log) error {
		for _, item := range sorted {
			fmt.Printf("%v\n", item)
		}
		return nil
	}))
```