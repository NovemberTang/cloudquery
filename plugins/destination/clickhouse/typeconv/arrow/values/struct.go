package values

import (
	"fmt"
	"strings"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"golang.org/x/exp/maps"
)

func buildStruct(builder *array.StructBuilder, value *map[string]any) error {
	if value == (*map[string]any)(nil) {
		builder.AppendNull()
		return nil
	}
	return appendStruct(builder, *value)
}

func appendStruct(builder *array.StructBuilder, value map[string]any) error {
	fields := builder.Type().(*arrow.StructType).Fields()
	remaining := make(map[string]struct{}, len(fields))
	for _, field := range fields {
		remaining[field.Name] = struct{}{}
	}

	for i, field := range fields {
		if err := buildValue(builder.FieldBuilder(i), value[field.Name]); err != nil {
			return err
		}
		delete(remaining, field.Name)
	}

	if len(remaining) > 0 {
		// ClickHouse SDK actually scans all tuple fields, even NULL values
		return fmt.Errorf("unresolved struct fields: [%s]", strings.Join(maps.Keys(remaining), ", "))
	}
	return nil
}
