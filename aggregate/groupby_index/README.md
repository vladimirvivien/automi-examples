## Aggregate - Group by Index Position

The example in this directory shows how to use Automi to Go slice data items
based on values at a specified index.

First the example sets up the data source as a slice of slices which will emit
each slice as a streamed item:

```go
source := sources.Slice([][]string{ ... }
```
Next, the code sets up the stream from the slice source and defines a batch
operator to collect all of the incoming slice items into a slice of slices
`[][]string`. Then, the GroupByIndex operator groups the slices based on the
value at index (3). The result is a map[string][][]string that contains the grouping.

```go
	// Create new stream from source
	strm := stream.From(source)

	// Setup stream operations
	strm.Run(
		// Batch each incoming `[]string` into --> [][]string
		window.Batch[[]string](),

		// Group batch items by value at index 3 into --> map[string][][]string
		exec.GroupByIndex[[][]string](3),
	)

	// Define sink operator to handle items
	strm.Into(sinks.Func(func(items map[string][][]string) error {
		for group, item := range items {
			fmt.Println("Group:", group)
			for _, slice := range item {
				fmt.Printf("%v\n", slice)
			}
		}
		return nil
	}))
```