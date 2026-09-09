package media

import (
	"fmt"
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
)

// outputFor builds a tool output schema explicitly, so that AddTool skips inference
// and uses it as is. It panics on failure: schemas are built once at startup from
// static types, so an error means the types cannot be described at all.
func outputFor[T any]() *jsonschema.Schema {
	schema, err := jsonschema.For[T](nil)
	if err != nil {
		panic(fmt.Errorf("build output schema for %s: %w", reflect.TypeFor[T](), err))
	}

	return schema
}
