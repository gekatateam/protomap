package protomap_test

import (
	"encoding/base64"
	"fmt"
)

const (
	testProto   = "./testdata/payload.proto"
	testBinary  = "./testdata/payload.binpb"
	testJson    = "./testdata/payload.json"
	testMessage = "protomap.test.Test"

	testIntersProto   = "./testdata/withtimeduration.proto"
	testIntersBinary  = "./testdata/withtimeduration.binpb"
	testIntersJson    = "./testdata/withtimeduration.json"
	testIntersMessage = "protomap.test.WithTimeDuration"
)

func setExpectedKeysWithTypes(in map[string]any) (map[string]any, error) {
	var err error
	in["Binary"], err = base64.StdEncoding.DecodeString(in["Binary"].(string))
	if err != nil {
		return nil, fmt.Errorf("binary decoding failed: %w", err)
	}

	in["Int"] = int64(in["Int"].(float64))
	in["Uint"] = uint64(in["Uint"].(float64))

	intsList := make([]any, 0)
	for _, v := range in["Inner"].(map[string]any)["List"].([]any) {
		intsList = append(intsList, int64(v.(float64)))
	}
	in["Inner"].(map[string]any)["List"] = intsList

	for k, v := range in["IntMap"].(map[string]any) {
		in["IntMap"].(map[string]any)[k] = int64(v.(float64))
	}

	return in, nil
}

func setInputKeysWithTypes(in map[string]any) (map[string]any, error) {
	var err error
	in["Binary"], err = base64.StdEncoding.DecodeString(in["Binary"].(string))
	if err != nil {
		return nil, fmt.Errorf("binary decoding failed: %w", err)
	}

	return in, nil
}
