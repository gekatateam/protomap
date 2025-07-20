package protomap_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/gekatateam/protomap"
)

func TestEncoder_EncodeToBinaryThenDecode(t *testing.T) {
	mapper, err := protomap.NewMapper(nil, testProto)
	if err != nil {
		t.Fatalf("decoder creation failed: %v", err)
	}

	binary, err := os.ReadFile(testBinary)
	if err != nil {
		t.Fatalf("binary data reading failed: %v", err)
	}

	tjson, err := os.ReadFile(testJson)
	if err != nil {
		t.Fatalf("json data reading failed: %v", err)
	}

	input := make(map[string]any)
	err = json.Unmarshal(tjson, &input)
	if err != nil {
		t.Fatalf("json data unmarshaling failed: %v", err)
	}

	input, err = setInputKeysWithTypes(input)
	if err != nil {
		t.Fatalf("map input preparation failed: %v", err)
	}

	result, err := mapper.Encode(input, testMessage)
	if err != nil {
		t.Fatalf("map input encoding failed: %v", err)
	}

	// need to find better way to validate result
	// we did this because fields order are random
	if len(binary) != len(result) {
		t.Log("------ expected -----")
		t.Logf("%v", binary)
		t.Log("------ result -----")
		t.Logf("%v", result)
		t.Fatal("expected and result length are not equal")
	}

	expected := make(map[string]any)
	err = json.Unmarshal(tjson, &expected)
	if err != nil {
		t.Fatalf("expected json unmarshaling failed: %v", err)
	}

	expected, err = setExpectedKeysWithTypes(expected)
	if err != nil {
		t.Fatalf("expected data preparation failed: %v", err)
	}

	decoderesult, err := mapper.Decode(binary, testMessage)
	if err != nil {
		t.Fatalf("binary data decoding failed: %v", err)
	}

	if !reflect.DeepEqual(expected, decoderesult) {
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
