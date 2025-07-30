package protomap_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/bufbuild/protocompile"
	"github.com/gekatateam/protomap"
	"github.com/gekatateam/protomap/interceptors"
)

func TestDecoder_DecodeToJson(t *testing.T) {
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

	expected := make(map[string]any)
	err = json.Unmarshal(tjson, &expected)
	if err != nil {
		t.Fatalf("json data unmarshaling failed: %v", err)
	}

	expected, err = setExpectedKeysWithTypes(expected)
	if err != nil {
		t.Fatalf("expected data preparation failed: %v", err)
	}

	result, err := mapper.Decode(binary, testMessage)
	if err != nil {
		t.Fatalf("binary data decoding failed: %v", err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Log("------ expected -----")
		for k, v := range expected {
			t.Logf("%v %v %T", k, v, v)
		}
		t.Log("------ result -----")
		for k, v := range result.(map[string]any) {
			t.Logf("%v %v %T", k, v, v)
		}
		t.Fatal("expected and result are not equal")
	}
}

func TestDecoder_DecodeWithInters(t *testing.T) {
	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{}),
	}

	mapper, err := protomap.NewMapper(&compiler, testIntersProto)
	if err != nil {
		t.Fatalf("decoder creation failed: %v", err)
	}

	binary, err := os.ReadFile(testIntersBinary)
	if err != nil {
		t.Fatalf("binary data reading failed: %v", err)
	}

	tjson, err := os.ReadFile(testIntersJson)
	if err != nil {
		t.Fatalf("json data reading failed: %v", err)
	}

	expected := make(map[string]any)
	err = json.Unmarshal(tjson, &expected)
	if err != nil {
		t.Fatalf("json data unmarshaling failed: %v", err)
	}

	expected["Ts"], err = time.Parse(time.RFC3339, expected["Ts"].(string))
	if err != nil {
		t.Fatalf("json time unmarshaling failed: %v", err)
	}

	expected["Dur"], err = time.ParseDuration(expected["Dur"].(string))
	if err != nil {
		t.Fatalf("json duration unmarshaling failed: %v", err)
	}

	result, err := mapper.Decode(binary, testIntersMessage, interceptors.DurationDecoder, interceptors.TimeDecoder)
	if err != nil {
		t.Fatalf("binary data decoding failed: %v", err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Log("------ expected -----")
		for k, v := range expected {
			t.Logf("%v %v %T", k, v, v)
		}
		t.Log("------ result -----")
		for k, v := range result.(map[string]any) {
			t.Logf("%v %v %T", k, v, v)
		}
		t.Fatal("expected and result are not equal")
	}
}
