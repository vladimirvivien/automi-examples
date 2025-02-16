# Automi Sink Examples

Sinks are Automi operators added at the end of a stream to collect streamed items. 
Automi supports several built-in sink implementations as showcased here.

* [csv](./csv) - Collects streamed items as CSV records.
* [discard](./discard/) - Discards streamed items into a sink
* [func](./func/) - Uses user-defined function as a stream sink
* [slice](./slice) - Collects streamed intems into a Go slice
* [slog](./slog/) - Uses Go `slog` Logger to collect stream items
* [writer](./writer) - Collects stream items using an `io.Writer` as sink