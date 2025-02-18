# Stream Custom Types

This example shows how to map to and stream data using 
custom types. The example loads data from a CSV and maps
each row to values of a specified Go type.

```go
type scientist struct {
	FirstName string
	LastName  string
	Title     string
	BornYear  int
}

func main() {
    source := sources.CSV(src)
    sink := sinks.CSV(snk)

    // stream from source
	stream := stream.From(source)

    // setup execution flow
	stream.Run(
		// map csv row to struct scientist
		exec.Map(func(ctx context.Context, cs []string) scientist {
			yr, _ := strconv.Atoi(cs[3])
			return scientist{
				FirstName: cs[1],
				LastName:  cs[0],
				Title:     cs[2],
				BornYear:  yr,
			}
		}),

		// apply data filter
		exec.Filter(func(ctx context.Context, cs scientist) bool {
			return (cs.BornYear > 1930)
		}),

		// remap value of type scientst to []string
		exec.Map(func(ctx context.Context, cs scientist) []string {
			return []string{cs.FirstName, cs.LastName, cs.Title}
		}),
	)
	// stream result into sink
	stream.Into(sink)
}
```