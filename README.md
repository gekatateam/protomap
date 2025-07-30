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

/* https://pkg.go.dev/github.com/bufbuild/protocompile#Compiler that will be used to compile .proto files; may be nil */
compiler := protocompile.Compiler{
    Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{}),
}

mapper, err := protomap.NewMapper(compiler, protofiles...)
if err != nil {
    panic(err)
}
```

### Decode your message from bytes slice to map
```go
/* full name of the message to decode */
messageName := "protomap.test.WithTimeDuration" 
result, err := mapper.Decode(binaryData, messageName)
if err != nil {
    panic(err)
}
```

### Encode your map to bytes slice
```go
/* full name of the message to encode */
messageName := "protomap.test.WithTimeDuration"
binaryData, err := mapper.Encode(gomap, messageName)
if err != nil {
    panic(err)
}
```

## Interceptors
Interceptors is a functions that allows you to encode/decode custom messages to Go types, for example, `time.Time` <-> `google.protobuf.Timestamp`. 

You need to pass interceptor(s) to `Encode`/`Decode` mapper methods to apply it:
```go
messageName := "protomap.test.WithTimeDuration"
binaryData, err := mapper.Encode(gomap, messageName, interceptors.DurationEncoder, interceptors.TimeEncoder)
if err != nil {
    panic(err)
}

A few ready functions you can find in [interceptors](interceptors/) dir.
