// Package jsontype provides a JSON wrapper type that round-trips correctly
// between Go, JSON, and SQL drivers (Postgres jsonb, MySQL json, SQLite TEXT).
//
// JSON is a thin alias over json.RawMessage that implements:
//
//   - sql.Scanner / driver.Valuer for database round-tripping
//   - json.Marshaler / json.Unmarshaler so it serialises as raw JSON
//     instead of the base64-encoded []byte the default would produce
//   - MarshalGQL / UnmarshalGQL for gqlgen scalar adapters (the methods
//     are duck-typed; jsontype itself does not import gqlgen)
package jsontype

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
)

// JSON is a json.RawMessage that round-trips through SQL drivers and gqlgen.
type JSON json.RawMessage

// Value implements database/sql/driver.Valuer.
//
// An empty JSON returns nil so it lands in the database as SQL NULL rather
// than the literal text "null".
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

// Scan implements database/sql.Scanner. It accepts []byte, string, or nil
// (mapped to the JSON null literal).
func (j *JSON) Scan(value any) error {
	if value == nil {
		*j = JSON("null")
		return nil
	}
	var b []byte
	switch v := value.(type) {
	case []byte:
		if len(v) > 0 {
			b = make([]byte, len(v))
			copy(b, v)
		}
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("jsontype: cannot scan %T into JSON", value)
	}
	*j = JSON(b)
	return nil
}

// MarshalJSON returns the raw JSON bytes without base64 encoding.
func (j JSON) MarshalJSON() ([]byte, error) {
	return json.RawMessage(j).MarshalJSON()
}

// UnmarshalJSON deserialises raw JSON bytes into JSON. Returns an error if
// the input is not valid JSON (json.RawMessage.UnmarshalJSON does not
// validate, so we wrap json.Valid for safer round-tripping).
func (j *JSON) UnmarshalJSON(b []byte) error {
	if !json.Valid(b) {
		return fmt.Errorf("jsontype: invalid JSON: %q", b)
	}
	*j = JSON(append([]byte(nil), b...))
	return nil
}

// MarshalGQL writes the JSON bytes to w. Implements the gqlgen Marshaler
// interface (duck-typed; this package does not import gqlgen).
func (j JSON) MarshalGQL(w io.Writer) {
	b, _ := json.Marshal(j)
	_, _ = w.Write(b)
}

// UnmarshalGQL parses an input from gqlgen — string, []byte, or
// map[string]any — into JSON.
func (j *JSON) UnmarshalGQL(input any) error {
	switch v := input.(type) {
	case string:
		return json.Unmarshal([]byte(v), j)
	case []byte:
		return json.Unmarshal(v, j)
	case map[string]any:
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return json.Unmarshal(b, j)
	default:
		return fmt.Errorf("jsontype: cannot unmarshal %T into JSON", input)
	}
}

// From marshals an arbitrary value into a *JSON. Convenient for passing
// arbitrary structures into a database column without manual marshalling.
func From(v any) (*JSON, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	out := JSON(b)
	return &out, nil
}
