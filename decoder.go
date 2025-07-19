package protomap

import (
	"context"
	"fmt"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type Decoder struct {
	files linker.Files
}

func NewDecoder(compiler *protocompile.Compiler, files ...string) (*Decoder, error) {
	if compiler == nil {
		compiler = &protocompile.Compiler{
			Resolver: &protocompile.SourceResolver{},
		}
	}

	f, err := compiler.Compile(context.Background(), files...)
	if err != nil {
		return nil, err
	}

	return &Decoder{files: f}, nil
}

func (d *Decoder) Unmarshal(data []byte, filepath, messageName string) (map[string]any, error) {
	f := d.files.FindFileByPath(filepath)
	if f == nil {
		return nil, ErrNoSuchFile
	}

	msgs := f.Messages()
	if msgs == nil {
		return nil, ErrNoMessages
	}

	desc := msgs.ByName(protoreflect.Name(messageName))
	if desc == nil {
		return nil, ErrNoSuchMessage
	}

	message := dynamicpb.NewMessage(desc)
	if err := proto.Unmarshal(data, message); err != nil {
		return nil, err
	}

	return MessageToMap(message)
}

func MessageToMap(message protoreflect.Message) (map[string]any, error) {
	fields := message.Descriptor().Fields()
	result := make(map[string]any, fields.Len())

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)

		if field.IsList() {
			list := message.Get(field).List()
			slice := make([]any, 0, list.Len())
			for j := 0; j < list.Len(); j++ {
				value, err := ProtoToGoValue(field.Kind(), list.Get(j))
				if err != nil {
					return nil, fmt.Errorf("%v.%v: %w", string(field.Name()), j, err)
				}
				slice = append(slice, value)
			}

			result[string(field.Name())] = slice
			continue
		}

		if field.IsMap() {
			pmap := message.Get(field).Map()
			gomap := make(map[string]any, pmap.Len())
			mapvaluekind := field.MapValue().Kind()

			var err error
			var failedKey string
			pmap.Range(func(mk protoreflect.MapKey, v protoreflect.Value) bool {
				value, convertErr := ProtoToGoValue(mapvaluekind, v)
				if convertErr != nil {
					err = convertErr
					failedKey = mk.String()
					return false
				}
				gomap[mk.String()] = value
				return true
			})

			if err != nil {
				return nil, fmt.Errorf("%v.%v: %w", string(field.Name()), failedKey, err)
			}

			result[string(field.Name())] = gomap
			continue
		}

		value, err := ProtoToGoValue(field.Kind(), message.Get(field))
		if err != nil {
			return nil, fmt.Errorf("%v: %w", string(field.Name()), err)
		}
		result[string(field.Name())] = value
	}

	return result, nil
}
