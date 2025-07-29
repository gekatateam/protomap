package interceptors

import (
	"math"
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func TimeDecoder(message protoreflect.Message) (result any, applied bool, err error) {
	if message.Descriptor().FullName() != "google.protobuf.Timestamp" {
		return nil, false, nil
	}

	seconds := message.Get(message.Descriptor().Fields().ByName("seconds")).Int()
	nanos := message.Get(message.Descriptor().Fields().ByName("nanos")).Int()

	return time.Unix(seconds, nanos).UTC(), true, nil
}

func DurationDecoder(message protoreflect.Message) (result any, applied bool, err error) {
	if message.Descriptor().FullName() != "google.protobuf.Duration" {
		return nil, false, nil
	}

	seconds := message.Get(message.Descriptor().Fields().ByName("seconds")).Int()
	nanos := message.Get(message.Descriptor().Fields().ByName("nanos")).Int()

	d := time.Duration(seconds) * time.Second
	overflow := d/time.Second != time.Duration(seconds)
	d += time.Duration(nanos) * time.Nanosecond
	overflow = overflow || (seconds < 0 && nanos < 0 && d > 0)
	overflow = overflow || (seconds > 0 && nanos > 0 && d < 0)
	if overflow {
		switch {
		case seconds < 0:
			return time.Duration(math.MinInt64), true, nil
		case seconds > 0:
			return time.Duration(math.MaxInt64), true, nil
		}
	}
	return d, true, nil
}
