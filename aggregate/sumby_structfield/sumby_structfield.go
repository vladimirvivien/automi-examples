package main

import (
    "context"
    "fmt"

    "github.com/vladimirvivien/automi/operators/exec"
    "github.com/vladimirvivien/automi/operators/window"
    "github.com/vladimirvivien/automi/sinks"
    "github.com/vladimirvivien/automi/sources"
    "github.com/vladimirvivien/automi/stream"
)

type stat struct {
    Zip                 string
    Count, Female, Male int
}

func main() {
    data := []stat{
        {Zip: "10452", Count: 17, Female: 12, Male: 5},
        {Zip: "10453", Count: 14, Female: 7, Male: 7},
        {Zip: "10454", Count: 18, Female: 8, Male: 10},
        {Zip: "10455", Count: 27, Female: 17, Male: 10},
        {Zip: "10456", Count: 5, Female: 3, Male: 2},
        {Zip: "10458", Count: 52, Female: 25, Male: 27},
        {Zip: "10459", Count: 7, Female: 5, Male: 2},
        {Zip: "10460", Count: 27, Female: 20, Male: 7},
        {Zip: "10461", Count: 49, Female: 26, Male: 23},
    }

    // Create stream from slice source
    strm := stream.From(sources.Slice(data)).Run(
        // Batch the data into a slice of stat
        window.Batch[stat](),
        // Sum by the specified struct field
        exec.SumByStructField[[]stat]("Male"),
    )

    // Setup sink to collect sum
    strm.Into(sinks.Func(func(val float64) error {
        fmt.Printf("Total male is %.0f\n", val)
        return nil
    }))

    // Open the stream with context
    if err := <-strm.Open(context.Background()); err != nil {
        fmt.Println(err)
        return
    }
}