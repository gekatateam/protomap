package interceptors

import (
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func TimeEncoder(input any, message protoreflect.Message) (applied bool, err error) {
	if message.Descriptor().FullName() != "google.protobuf.Timestamp" {
		return false, nil
	}

	t, ok := input.(time.Time)
	if !ok {
		return false, nil
	}

	message.Set(message.Descriptor().Fields().ByName("seconds"), protoreflect.ValueOfInt64(t.Unix()))
	message.Set(message.Descriptor().Fields().ByName("nanos"), protoreflect.ValueOfInt32(int32(t.Nanosecond())))

	return true, nil
}

func DurationEncoder(input any, message protoreflect.Message) (applied bool, err error) {
	if message.Descriptor().FullName() != "google.protobuf.Duration" {
		return false, nil
	}

	d, ok := input.(time.Duration)
	if !ok {
		return false, nil
	}

	nanos := d.Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9

	message.Set(message.Descriptor().Fields().ByName("seconds"), protoreflect.ValueOfInt64(secs))
	message.Set(message.Descriptor().Fields().ByName("nanos"), protoreflect.ValueOfInt32(int32(nanos)))

	return true, nil
}
