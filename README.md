# automi-examples

Repository for project Automi code examples. This repository provides a variety of examples demonstrating how to use Automi for different use cases.

## Aggregate Examples
The [aggregate](./aggregate) directory contains examples of Automi aggregation operators and how they are used.

*   [groupby\_index](./aggregate/groupby_index/) - Batches and groups slice/array values based on index position.
*   [groupby\_mapkey](./aggregate/groupby_mapkey) - Batches and groups streaming map values based on their keys.
*   [groupby\_structfield](./aggregate/groupby_structfield) - Batches and groups struct values based on field names.
*   [sort\_slice](./aggregate/sort_slice) - Sorts batched items using the natural sort sequence of streamed items.
*   [sort\_withfunc](./aggregate/sort_withfunc/) - Shows how to sort items using a custom sorting function.
*   [sortby\_index](./aggregate/sortby_index/) - Shows how to sort streaming slice values based on selected index.
*   [sortby\_mapkey](./aggregate/sortby_mapkey) - Shows how to sort streaming map values based on their keys.
*   [sortby\_structfield](./aggregate/sortby_structfield) - Shows how to sort streaming struct values based on field names.
*   [sum](./aggregate/sum) - Shows how to batch and add streaming numeric values.
*   [sumby\_index](./aggregate/sumby_index/) - Shows how to add straming slice values based on selected index.
*   [sumby\_mapkey](./aggregate/sumby_mapkey) - Shows how to add streaming map values based on their keys.

## Other Examples

*   [Custom type](./customtype) - Shows how to stream values of custom types.
*   [Error handling](./error) - Examples showing how to set up error handling.
*   [gRPC streaming](./grpc) - Examples using Automi to stream from gRPC streaming services.
*   [Logging](./logging) - Examples of logging stream events at runtime.
*   [MD5](./md5) - Implementation of the MD5 example from Sameer Ajmani's blog on [Concurrency Pattern](https://blog.golang.org/pipelines).
*   [Network](./net) - Examples showing streaming data from TCP sockets and HTTP requests.
    *   [http](./net/http) - Demonstrates how to create an HTTP server that uses Automi to process incoming request data and return a response to the client.
    *   [socket](./net/socket) - Demonstrates how to stream data from TCP sockets.
*   [Sinks](./sinks) - Examples of built-in sink components.
*   [Sources](./sources) - Examples of built-in source components.
    *   [CSV](./sources/csv) - Demonstrates Automi's support for streaming data from CSV sources.
    *   [Scanner](./sources/scanner) - Shows how to use Automi to stream data from a Go `bufio.Scanner` source.
    *   [Slice0](./sources/slice0) - Shows how to stream data from a Go slice source.
*   [Wordcount](./wordcount) - A simple wordcount example using Automi operators.
