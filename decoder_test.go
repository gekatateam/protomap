package protomap_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/gekatateam/protomap"
)

const (
	testProto   = "./testdata/payload.proto"
	testBinary  = "./testdata/payload.binpb"
	testJson    = "./testdata/payload.json"
	testMessage = "Test"
)

func TestDecoder_E2E(t *testing.T) {
	decoder, err := protomap.NewDecoder(nil, testProto)
	if err != nil {
		t.Fatalf("decoder create failed: %v", err)
	}

	binary, err := os.ReadFile(testBinary)
	if err != nil {
		t.Fatalf("binary data read failed: %v", err)
	}

	tjson, err := os.ReadFile(testJson)
	if err != nil {
		t.Fatalf("json data read failed: %v", err)
	}

	expected := make(map[string]any)
	err = json.Unmarshal(tjson, &expected)
	if err != nil {
		t.Fatalf("json data unmarshal failed: %v", err)
	}

	expected, err = setExpectedKeysWithTypes(expected)
	if err != nil {
		t.Fatalf("expected data preparation failed: %v", err)
	}

	result, err := decoder.Unmarshal(binary, testProto, testMessage)
	if err != nil {
		t.Fatalf("binary data unmarshal failed: %v", err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Log("------ expected -----")
		for k, v := range expected {
			t.Logf("%v %v %T", k, v, v)
		}
		t.Log("------ result -----")
		for k, v := range result {
			t.Logf("%v %v %T", k, v, v)
		}
		t.Fatal("expected and result are not equal")
	}
}

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
