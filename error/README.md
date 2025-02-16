## Using `StreamResult` for Stream Error Handling


This example (builds on md5) shows how to use type `api.StreamResult` to return
a value or an error. The type also allows you to specify how to handle the error
in the stream.

```go
exec.Map(func(ctx context.Context, info walkInfo) api.StreamResult {
	if strings.HasPrefix(filepath.Base(info.path), "a") {
		return api.StreamResult{
			Err:    errors.New("encountered file that starts with letter `a`: skipping"),
			Action: api.ActionSkipItem,
		}
	}
	return api.StreamResult{Value: info.path}
}),
```
