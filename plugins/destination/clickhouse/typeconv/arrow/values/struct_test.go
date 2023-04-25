package values

import (
	"testing"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/cloudquery/plugin-sdk/v2/types"
	"github.com/google/uuid"
)

func Test_struct(t *testing.T) {
	_true := true
	for _, tc := range []testCase{
		{
			_type:    arrow.StructOf(arrow.Field{Name: "f1", Type: new(arrow.BooleanType)}),
			value:    map[string]any{"f1": true},
			expected: map[string]any{"f1": true},
		},
		{
			_type:    arrow.StructOf(arrow.Field{Name: "f1", Type: new(arrow.BooleanType), Nullable: true}),
			value:    map[string]any{"f1": &_true},
			expected: map[string]any{"f1": true},
		},
		{
			_type:    arrow.StructOf(arrow.Field{Name: "f1_null", Type: new(arrow.BooleanType), Nullable: true}),
			value:    map[string]any{"f1_null": nil},
			expected: map[string]any{"f1_null": nil},
		},
		{
			_type: arrow.StructOf(
				arrow.Field{
					Name: "bool_list",
					Type: arrow.ListOfNonNullable(new(arrow.BooleanType)),
				},
				arrow.Field{
					Name:     "bool_list_nullable",
					Type:     arrow.ListOf(new(arrow.BooleanType)),
					Nullable: true,
				},
			),
			value:    map[string]any{"bool_list": []bool{true}, "bool_list_nullable": nil},
			expected: map[string]any{"bool_list": marshalList(t, []bool{true}), "bool_list_nullable": nil},
		},
		{
			_type: arrow.StructOf(
				arrow.Field{
					Name:     "uuid_list",
					Type:     arrow.ListOf(types.NewUUIDType()),
					Nullable: true,
				},
				arrow.Field{
					Name: "struct",
					Type: arrow.StructOf(
						arrow.Field{
							Name: "bool_list",
							Type: arrow.ListOfNonNullable(new(arrow.BooleanType)),
						},
						arrow.Field{
							Name:     "bool_list_nullable",
							Type:     arrow.ListOf(new(arrow.BooleanType)),
							Nullable: true,
						},
					),
				},
			),
			value: map[string]any{
				"uuid_list": []*uuid.UUID{&uuid.NameSpaceOID},
				"struct":    map[string]any{"bool_list": []bool{true}, "bool_list_nullable": nil},
			},
			expected: map[string]any{
				"uuid_list": marshalList(t, []uuid.UUID{uuid.NameSpaceOID}),
				"struct":    map[string]any{"bool_list": marshalList(t, []bool{true}), "bool_list_nullable": nil},
			},
		},
	} {
		ensureRecord(t, tc)
	}
}
