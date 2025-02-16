## CSV Source

This example demonstrates Automi's support for streaming data from
CSV sources. First, the code opens the file with the CSV data and
configures the Automi CSV source:

```go
srcFile, err := os.Open("./stats.csv")
csv := sources.CSV(srcFile).
	CommentChar('#'). // sets comment charcter, default is '#'
	DelimChar(',').   // sets delimiter char, default is  ','
	HasHeaders()      // indicate CSV has headers
```

The remainder of the code streams the data from the CSV source
and processes each row using a `Map` and `Filter` operators:

```go
// Create new stream from CSV source
strm := stream.From(csv)

strm.Run(
	// select first 6 cols per row
	exec.Map(func(_ context.Context, row []string) []string {
		return row[:6]
	}),

	// filter out rows with zero participants
	exec.Filter(func(_ context.Context, row []string) bool {
		count, err := strconv.Atoi(row[1])
		if err != nil {
			count = 0
		}
		return (count > 0)
	}),
)

// sink result in a collector function to prints it
strm.Into(sinks.Func(func(row []string) error {
	fmt.Println(row)
	return nil
}))
```