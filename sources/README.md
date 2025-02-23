# Automi Sources Examples

Automi sources are operators that are placed at the start of a stream as the source of data. 
The source can emit data from other sources or produce the data itself.

The following shows examples for supported Automi emitters.

* [channel](./channel) - Using Go channel as stream source.
* [csv](./csv) - Shows how to use records from CSV as stream source
* [reader](./reader) - Uses a Go  `io.Reader` to emit stream items
* [scanner](./scanner) - Uses buffered IO to stream items from an `io.Reader`
* [slice0](./slice0) - Eemit items of built-in types from a slice source
* [slice1](./slice1) - Emit items of custom type from a slice source