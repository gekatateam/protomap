# protomap

Protomap is a `map[string]any` <-> `protobuf binary` encoder/decoder based on [protoreflect](https://pkg.go.dev/google.golang.org/protobuf/reflect/protoreflect).

Like json `Unmarshal`/`Marshal` to/from generic map, but for protobuf.


## Why?

To operate protobuf-encoded messages without updating generated code each time when your service need to work with tons of updated/new `.proto`.


## How?

`Mapper` creates new message from descriptor every `Encode`/`Decode` call, so, it should be concurrent-safe to encode/decode multiple messages at the same time using one `Mapper`.

### Create `Mapper`
```go
protofiles := []string{/* list of .proto files, which contains messages encode/decode to */}
compiler := /* https://pkg.go.dev/github.com/bufbuild/protocompile#Compiler that will be used to compile .proto files; may be nil */
mapper, err := protomap.NewMapper(compiler, protofiles...)
if err != nil {
    panic(err)
}
```

### Decode your message from bytes slice to map
```go
messageName := /* full name of the message to decode */
gomap, err := mapper.Decode(binaryData, messageName)
if err != nil {
    panic(err)
}
```

### Encode your map to bytes slice
```go
messageName := /* full name of the message to encode */
binaryData, err := mapper.Encode(gomap, messageName)
if err != nil {
    panic(err)
}
```
