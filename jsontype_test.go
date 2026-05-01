package jsontype

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestValue_EmptyReturnsNil(t *testing.T) {
	var j JSON
	v, err := j.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	if v != nil {
		t.Errorf("empty JSON: got %v, want nil", v)
	}
}

func TestValue_NonEmptyReturnsString(t *testing.T) {
	j := JSON(`{"a":1}`)
	v, err := j.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	if s, ok := v.(string); !ok || s != `{"a":1}` {
		t.Errorf("Value: got %T %v, want string %q", v, v, `{"a":1}`)
	}
}

func TestScan_Bytes(t *testing.T) {
	var j JSON
	if err := j.Scan([]byte(`{"x":2}`)); err != nil {
		t.Fatalf("Scan bytes: %v", err)
	}
	if string(j) != `{"x":2}` {
		t.Errorf("got %q, want %q", string(j), `{"x":2}`)
	}
}

func TestScan_String(t *testing.T) {
	var j JSON
	if err := j.Scan(`{"y":3}`); err != nil {
		t.Fatalf("Scan string: %v", err)
	}
	if string(j) != `{"y":3}` {
		t.Errorf("got %q, want %q", string(j), `{"y":3}`)
	}
}

func TestScan_Nil_BecomesJSONNull(t *testing.T) {
	var j JSON
	if err := j.Scan(nil); err != nil {
		t.Fatalf("Scan nil: %v", err)
	}
	if string(j) != "null" {
		t.Errorf("nil scan: got %q, want %q", string(j), "null")
	}
}

func TestScan_UnsupportedType(t *testing.T) {
	var j JSON
	if err := j.Scan(42); err == nil {
		t.Errorf("scan int: expected error")
	}
}

func TestMarshalJSON_RoundTrip(t *testing.T) {
	src := JSON(`{"k":"v","n":42}`)
	b, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(b) != `{"k":"v","n":42}` {
		t.Errorf("marshal: got %s, want raw JSON object", b)
	}

	var got JSON
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if string(got) != string(src) {
		t.Errorf("roundtrip: got %q, want %q", string(got), string(src))
	}
}

func TestUnmarshalJSON_InvalidJSON(t *testing.T) {
	var j JSON
	if err := j.UnmarshalJSON([]byte(`not json}`)); err == nil {
		t.Errorf("expected error on invalid JSON")
	}
}

func TestMarshalGQL_WritesRawJSON(t *testing.T) {
	j := JSON(`{"a":1}`)
	var buf bytes.Buffer
	j.MarshalGQL(&buf)
	if buf.Len() == 0 {
		t.Errorf("MarshalGQL wrote nothing")
	}
}

func TestUnmarshalGQL_String(t *testing.T) {
	var j JSON
	if err := j.UnmarshalGQL(`{"a":1}`); err != nil {
		t.Fatalf("UnmarshalGQL string: %v", err)
	}
	if string(j) != `{"a":1}` {
		t.Errorf("got %q", string(j))
	}
}

func TestUnmarshalGQL_Bytes(t *testing.T) {
	var j JSON
	if err := j.UnmarshalGQL([]byte(`{"a":1}`)); err != nil {
		t.Fatalf("UnmarshalGQL bytes: %v", err)
	}
}

func TestUnmarshalGQL_Map(t *testing.T) {
	var j JSON
	in := map[string]any{"a": 1, "b": "two"}
	if err := j.UnmarshalGQL(in); err != nil {
		t.Fatalf("UnmarshalGQL map: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(j, &got); err != nil {
		t.Fatalf("re-unmarshal: %v", err)
	}
	if !reflect.DeepEqual(got["a"], float64(1)) || got["b"] != "two" {
		t.Errorf("map round-trip wrong: %v", got)
	}
}

func TestUnmarshalGQL_UnsupportedType(t *testing.T) {
	var j JSON
	if err := j.UnmarshalGQL(123); err == nil {
		t.Errorf("expected error on int input")
	}
}

func TestFrom_ArbitraryStruct(t *testing.T) {
	type S struct {
		A int    `json:"a"`
		B string `json:"b"`
	}
	j, err := From(S{A: 1, B: "x"})
	if err != nil {
		t.Fatalf("From: %v", err)
	}
	if string(*j) != `{"a":1,"b":"x"}` {
		t.Errorf("got %q", string(*j))
	}
}

func TestFrom_UnmarshallableInput(t *testing.T) {
	if _, err := From(make(chan int)); err == nil {
		t.Errorf("expected error for unmarshalable input")
	}
}
