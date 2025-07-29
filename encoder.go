package protomap

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func (e *Mapper) Encode(data any, messageName string) ([]byte, error) {
	desc, err := e.r.FindMessageByName(protoreflect.FullName(messageName))
	if err != nil {
		return nil, err
	}

	if desc == nil {
		return nil, ErrNoSuchMessage
	}

	message := dynamicpb.NewMessage(desc.Descriptor())
	if err := AnyToMessage(data, message); err != nil {
		return nil, err
	}

	return proto.Marshal(message)
}
