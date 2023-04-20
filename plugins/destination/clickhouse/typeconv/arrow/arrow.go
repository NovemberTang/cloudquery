package arrow

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/column"
	"github.com/ClickHouse/clickhouse-go/v2/lib/timezone"
	"github.com/apache/arrow/go/v12/arrow"
	"github.com/cloudquery/cloudquery/plugins/destination/clickhouse/util"
	"github.com/cloudquery/plugin-sdk/v2/types"
)

func fieldFromColumn(col column.Interface) (*arrow.Field, error) {
	fieldName := util.UnquoteID(col.Name())
	switch col := col.(type) {
	case *column.Bool:
		return &arrow.Field{Name: fieldName, Type: new(arrow.BooleanType)}, nil

	case *column.UInt8:
		return &arrow.Field{Name: fieldName, Type: new(arrow.Uint8Type)}, nil
	case *column.UInt16:
		return &arrow.Field{Name: fieldName, Type: new(arrow.Uint16Type)}, nil
	case *column.UInt32:
		return &arrow.Field{Name: fieldName, Type: new(arrow.Uint32Type)}, nil
	case *column.UInt64:
		return &arrow.Field{Name: fieldName, Type: new(arrow.Uint64Type)}, nil
	case *column.Int8:
		return &arrow.Field{Name: fieldName, Type: new(arrow.Int8Type)}, nil
	case *column.Int16:
		return &arrow.Field{Name: fieldName, Type: new(arrow.Int16Type)}, nil
	case *column.Int32:
		return &arrow.Field{Name: fieldName, Type: new(arrow.Int32Type)}, nil
	case *column.Int64:
		return &arrow.Field{Name: fieldName, Type: new(arrow.Int64Type)}, nil

	case *column.Float32:
		return &arrow.Field{Name: fieldName, Type: new(arrow.Float32Type)}, nil
	case *column.Float64:
		return &arrow.Field{Name: fieldName, Type: new(arrow.Float64Type)}, nil

	case *column.String:
		return &arrow.Field{Name: fieldName, Type: new(arrow.StringType)}, nil

	case *column.FixedString:
		// sadly, we need to parse manually here
		var byteWidth int
		if _, err := fmt.Sscanf(string(col.Type()), "FixedString(%d)", &byteWidth); err != nil {
			return nil, err
		}
		return &arrow.Field{Name: fieldName, Type: &arrow.FixedSizeBinaryType{ByteWidth: byteWidth}}, nil

	case *column.Date32:
		return &arrow.Field{Name: fieldName, Type: new(arrow.Date32Type)}, nil

	case *column.DateTime:
		// need to parse
		param := params(col.Type())
		name := strings.Trim(strings.TrimSpace(param), "'")
		tz, err := timezone.Load(name)
		if err != nil {
			return nil, err
		}

		return &arrow.Field{Name: fieldName, Type: &arrow.TimestampType{Unit: arrow.Second, TimeZone: tz.String()}}, nil

	case *column.DateTime64:
		// need to parse
		params := strings.Split(params(col.Type()), ",")
		var tz string
		precision, err := strconv.Atoi(params[0])
		if err != nil {
			return nil, err
		}
		var unit arrow.TimeUnit
		switch precision {
		case 0:
			unit = arrow.Second
		case 3:
			// This is the same as arrow.DATE64, so we need to canonize the schema
			unit = arrow.Millisecond
		case 6:
			unit = arrow.Microsecond
		case 9:
			unit = arrow.Nanosecond
		default:
			return nil, fmt.Errorf("unsupported DateTime64 precision: %d (supported values: 0,3,6,9)", precision)
		}

		if len(params) > 1 {
			name := strings.Trim(strings.TrimSpace(params[1]), "'")
			zone, err := timezone.Load(name)
			if err != nil {
				return nil, err
			}
			tz = zone.String()
		}

		return &arrow.Field{Name: fieldName, Type: &arrow.TimestampType{Unit: unit, TimeZone: tz}}, nil

	case *column.Decimal:
		var decimal arrow.DecimalType
		if precision := col.Precision(); precision <= 38 {
			decimal = &arrow.Decimal128Type{Precision: int32(precision), Scale: int32(col.Scale())}
		} else {
			decimal = &arrow.Decimal256Type{Precision: int32(precision), Scale: int32(col.Scale())}
		}
		return &arrow.Field{Name: fieldName, Type: decimal}, nil

	case *column.Array:
		base, err := fieldFromColumn(col.Base())
		if err != nil {
			return nil, err
		}
		return &arrow.Field{
			Name: fieldName,
			Type: arrow.ListOfField(*base),
		}, nil

	case *column.Nullable:
		base, err := fieldFromColumn(col.Base())
		if err != nil {
			return nil, err
		}
		return &arrow.Field{
			Name:     fieldName,
			Type:     base.Type,
			Nullable: true,
		}, nil

	case *column.Map:
		dataType, err := mapType(col)
		if err != nil {
			return nil, err
		}
		return &arrow.Field{Name: fieldName, Type: dataType}, nil

	case *column.Tuple:
		dataType, err := structType(col)
		if err != nil {
			return nil, err
		}
		return &arrow.Field{Name: fieldName, Type: dataType}, nil

	case *column.UUID:
		return &arrow.Field{Name: fieldName, Type: new(types.UUIDType)}, nil

	default:
		return &arrow.Field{Name: fieldName, Type: new(arrow.StringType)}, nil
	}
}

func Field(name, typ string) (*arrow.Field, error) {
	col, err := column.Type(typ).Column(name, time.UTC)
	if err != nil {
		return nil, err
	}

	return fieldFromColumn(col)
}
