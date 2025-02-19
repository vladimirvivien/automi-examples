Calculating MD5
===============

This examples calculates MD5 hash of files on a given path. The code
uses `emitPathFor` function to walk the path and emits file information on a channel that is then
used as a source for an Automi stream.

```go
stream := stream.From(sources.Chan(emitPathsFor(rootPath)))
```

Then the example sets up a set stream processing operation to calculate md5:

```go
exec.Map(func(ctx context.Context, filePath string) [3]any {
	data, err := os.ReadFile(filePath)
	sum := md5.Sum(data)
	return [3]any{filePath, sum, err}
}),
```
 
 Lastly, the result is printed in the sink:

 ```go
stream.Into(sinks.Func(func(items [3]any) error {
	file := items[0].(string)
	md5Sum := items[1].([md5.Size]byte)
	fmt.Printf("file %-64s md5 (%-16x)\n", file, md5Sum)
	return nil
}))
```

```bash
go run . -p ../
file ../.git/COMMIT_EDITMSG                                           md5 (aab2daf0bc3aad2b406b0bbfbcf8bf1e)
file ../.git/FETCH_HEAD                                               md5 (83bcb90bb00acd4a5c7b2c3e8eb6ef58)
file ../.git/HEAD                                                     md5 (cf7dd3ce51958c5f13fece957cc417fb)
file ../.git/ORIG_HEAD                                                md5 (f3e08307d05dd488a3f0d3e7dabfd510)
file ../.git/REBASE_HEAD                                              md5 (f3e08307d05dd488a3f0d3e7dabfd510)
file ../.git/config                                                   md5 (47d50efa537803aac019b0657e33eeca)
file ../.git/description                                              md5 (a0a7c3fff21f2aea3cfa1d0316dd816c)
```