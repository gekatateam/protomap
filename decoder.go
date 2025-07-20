package protomap

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func (d *Mapper) Decode(data []byte, messageName string) (map[string]any, error) {
	desc, err := d.r.FindMessageByName(protoreflect.FullName(messageName))
	if err != nil {
		return nil, err
	}

	if desc == nil {
		return nil, ErrNoSuchMessage
	}

	message := dynamicpb.NewMessage(desc.Descriptor())
	if err := proto.Unmarshal(data, message); err != nil {
		return nil, err
	}

	return MessageToMap(message)
}
