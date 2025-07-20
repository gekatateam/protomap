package protomap

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

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

func ProtoToGoValue(kind protoreflect.Kind, value protoreflect.Value) (any, error) {
	switch kind {
	case protoreflect.BoolKind:
		return value.Bool(), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return value.Int(), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return value.Uint(), nil
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return value.Float(), nil
	case protoreflect.StringKind:
		return value.String(), nil
	case protoreflect.BytesKind:
		return value.Bytes(), nil
	case protoreflect.EnumKind:
		return value.Enum(), nil
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return MessageToMap(value.Message())
	default:
		return nil, fmt.Errorf("unsupported field type: %s", kind)
	}
}

func MapToMessage(data map[string]any, message protoreflect.Message) error {
	fields := message.Descriptor().Fields()

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)

		value, ok := data[string(field.Name())]
		if !ok {
			if field.Cardinality() == protoreflect.Optional {
				continue
			}
			return fmt.Errorf("%v is not optional, but input data has no such key", field.Name())
		}

		if field.IsList() {
			slice, ok := value.([]any)
			if !ok {
				return fmt.Errorf("%v is a list, but input data field is not a slice", field.Name())
			}

			elemkind := field.Kind()
			protolist := message.NewField(field).List()
			for i, v := range slice {
				protovalue, err := GoValueToProto(field, elemkind, v)
				if err != nil {
					return fmt.Errorf("%v.%v: %w", field.Name(), i, err)
				}
				protolist.Append(protovalue)
			}
			message.Set(field, protoreflect.ValueOfList(protolist))
			continue
		}

		if field.IsMap() {
			gomap, ok := value.(map[string]any)
			if !ok {
				return fmt.Errorf("%v is a map, but input data field is not a map", field.Name())
			}

			keykind := field.MapKey().Kind()
			valkind := field.MapValue().Kind()
			protomap := message.NewField(field).Map()
			for k, v := range gomap {
				protokey, err := GoValueToProto(field, keykind, k)
				if err != nil {
					return fmt.Errorf("%v.%v key: %w", field.Name(), k, err)
				}

				protovalue, err := GoValueToProto(field, valkind, v)
				if err != nil {
					return fmt.Errorf("%v.%v key: %w", field.Name(), k, err)
				}

				protomap.Set(protokey.MapKey(), protovalue)
			}
			message.Set(field, protoreflect.ValueOfMap(protomap))
			continue
		}

		protovalue, err := GoValueToProto(field, field.Kind(), value)
		if err != nil {
			return fmt.Errorf("%v: %w", field.Name(), err)
		}

		message.Set(field, protovalue)
	}

	return nil
}

func GoValueToProto(desc protoreflect.FieldDescriptor, kind protoreflect.Kind, value any) (protoreflect.Value, error) {
	switch kind {
	case protoreflect.StringKind:
		v, err := AnyToString(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfString(v), nil
	case protoreflect.BoolKind:
		v, err := AnyToBoolean(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfBool(v), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		v, err := AnyToInteger(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfInt32(int32(v)), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		v, err := AnyToInteger(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfInt64(v), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		v, err := AnyToUnsigned(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfUint32(uint32(v)), nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		v, err := AnyToUnsigned(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfUint64(v), nil
	case protoreflect.FloatKind:
		v, err := AnyToFloat(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfFloat32(float32(v)), nil
	case protoreflect.DoubleKind:
		v, err := AnyToFloat(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfFloat64(v), nil
	case protoreflect.BytesKind:
		v, err := AnyToBytes(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfBytes(v), nil
	case protoreflect.EnumKind:
		v, err := AnyToInteger(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfEnum(protoreflect.EnumNumber(int32(v))), nil
	case protoreflect.MessageKind, protoreflect.GroupKind:
		v, ok := value.(map[string]any)
		if !ok {
			return protoreflect.Value{}, errors.New("field kind is message or gorup, but value not a map")
		}

		msg := dynamicpb.NewMessage(desc.Message())
		if err := MapToMessage(v, msg); err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfMessage(msg), nil
	default:
		return protoreflect.Value{}, fmt.Errorf("unsupported field type: %s", kind)
	}
}
