# jsontype

A `JSON` type that round-trips correctly between Go, JSON, and SQL drivers — Postgres `jsonb`, MySQL `json`, SQLite `TEXT`, anything that scans to `[]byte` or `string`.

Zero third-party dependencies.

## Install

```sh
go get github.com/ubgo/jsontype
```

## Quick start (with [Ent](https://entgo.io/) or any `database/sql` ORM)

```go
import "github.com/ubgo/jsontype"

type User struct {
    ID       int             `json:"id"`
    Settings jsontype.JSON   `json:"settings"`
}

// Settings is stored as jsonb / json / TEXT in the DB and scans round-trip
// via sql.Scanner / driver.Valuer.
```

## Quick start (with [gqlgen](https://gqlgen.com/) for GraphQL custom scalars)

`jsontype.JSON` already implements `MarshalGQL` / `UnmarshalGQL`. Add a scalar binding in your `gqlgen.yml`:

```yaml
models:
  JSON:
    model:
      - github.com/ubgo/jsontype.JSON
```

Then declare `scalar JSON` in your schema and use it on any field.

## API

| Method | Purpose |
|--------|---------|
| `Value() (driver.Value, error)` | `database/sql/driver.Valuer`. Empty JSON → SQL NULL. |
| `Scan(value any) error` | `database/sql.Scanner`. Accepts `[]byte`, `string`, `nil` (→ JSON `null` literal). |
| `MarshalJSON() ([]byte, error)` | Returns raw JSON bytes (no base64 wrapping). |
| `UnmarshalJSON(b []byte) error` | Validates input via `json.Valid` before storing. |
| `MarshalGQL(w io.Writer)` | gqlgen Marshaler. Duck-typed; jsontype does not import gqlgen. |
| `UnmarshalGQL(input any) error` | Accepts `string`, `[]byte`, or `map[string]any`. |
| `From(v any) (*JSON, error)` | Marshal an arbitrary value into a `*JSON`. |

## Why not `json.RawMessage` directly?

`json.RawMessage` doesn't implement `sql.Scanner` / `driver.Valuer`, and its `UnmarshalJSON` doesn't validate input — so `db.Scan(...)` and `json.Unmarshal(invalidBytes, &raw)` both fail in surprising ways. `jsontype.JSON` covers both gaps without leaving the standard library.

## License

Apache License 2.0. See [`LICENSE`](./LICENSE).
