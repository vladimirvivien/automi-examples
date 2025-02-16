## Word Count Example

This example demonstrates how to use Automi to perform word count on a text file. It reads data from a file, normalizes spaces, removes non-alphanumeric characters, and then counts the occurrences of each word.

### Code Overview

The example sets up a stream that reads data from a text file (`twotw.txt`), processes it to count word occurrences, and prints the results to the console.

```go
func main() {
	// Regular expressions for normalizing spaces and removing non-alphanumeric characters
	regSpaces := regexp.MustCompile(`\s+`)
	regNonAlpha := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	space := " "

	// Setup stream data source from file
	file, err := os.Open("./twotw.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
    defer file.Close()

	strm := stream.From(sources.Scanner(file, bufio.ScanLines))
	strm.Run(
		// Normalize all spaces to single space, then return as slice of words
		exec.Execute(func(_ context.Context, line []byte) []string {
			return strings.Split(regSpaces.ReplaceAllLiteralString(string(line), space), space)
		}),

		// Restreams the slice of words into individual words
		restream.Slice[[]string](),

		// Map each word to an occurrence tuple of Pair{word, 1}
		exec.Map(func(_ context.Context, word string) tuple.Pair[string, int] {
			word = regNonAlpha.ReplaceAllLiteralString(word, "")
			return tuple.Pair[string, int]{Val1: word, Val2: 1}
		}),

		// Batch the occurrence pairs into slices of []Pair{word, 1}
		window.Batch[tuple.Pair[string, int]](),

		// Group the occurrences by word --> map[word][]Pair{word, 1}:
		exec.GroupByStructField[[]tuple.Pair[string, int]]("Val1"),

		// Count the occurrences of each word from the group, return as []Pair{word, total}
		exec.Execute(func(_ context.Context, group map[any][]tuple.Pair[string, int]) []tuple.Pair[string, int] {
			var result []tuple.Pair[string, int]
			for key, pairs := range group {
				sum := 0
				for _, pair := range pairs {
					sum += pair.Val2
				}
				result = append(result, tuple.Pair[string, int]{Val1: key.(string), Val2: sum})
			}

			// Sort result
			slices.SortFunc(result, func(p1, p2 tuple.Pair[string, int]) int {
				return cmp.Compare(p1.Val1, p2.Val1)
			})

			return result
		}),
	)
}
```

The stream operations include:

1.  **Normalization and Splitting**: Normalizes spaces and splits lines into words using regular expressions.
2.  **Restreaming**: Converts the slice of words into individual words for processing.
3.  **Mapping**: Maps each word to an occurrence tuple (`Pair{word, 1}`).
4.  **Batching**: Batches the occurrence pairs into slices.
5.  **Grouping**: Groups the occurrences by word.
6.  **Counting**: Counts the occurrences of each word from the group.
7.  **Sorting**: Sorts the result by word.

Finally, the code sets up a sink to process and display the word counts:

```go
// Route the result to a sink that prints the word count
strm.Into(sinks.Func(func(wordp []tuple.Pair[string, int]) error {
	for _, pair := range wordp {
		fmt.Printf("%s: %d\n", pair.Val1, pair.Val2)
	}
	return nil
}))
```