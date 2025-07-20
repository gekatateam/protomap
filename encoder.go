package protomap

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func (e *Mapper) Encode(data map[string]any, filepath, messageName string) ([]byte, error) {
	f := e.files.FindFileByPath(filepath)
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
	if err := MapToMessage(data, message); err != nil {
		return nil, err
	}

	return proto.Marshal(message)
}
