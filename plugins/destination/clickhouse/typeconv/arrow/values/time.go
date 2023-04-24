package values

import (
	"fmt"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
)

func buildDate32Values(builder primitiveBuilder[arrow.Date32], value *time.Time) {
	switch {
	case value == nil, value == (*time.Time)(nil):
		builder.AppendNull()
	default:
		builder.Append(arrow.Date32FromTime(*value))
	}
}
func buildDate64Values(builder primitiveBuilder[arrow.Date64], value *time.Time) {
	switch {
	case value == nil, value == (*time.Time)(nil):
		builder.AppendNull()
	default:
		builder.Append(arrow.Date64FromTime(*value))
	}
}

func buildTimestampValues(builder *array.TimestampBuilder, value *time.Time) error {
	switch {
	case value == nil, value == (*time.Time)(nil):
		builder.AppendNull()
		return nil
	default:
		t, err := timeToTimestamp(*value, builder.Type().(*arrow.TimestampType))
		if err != nil {
			builder.AppendNull()
			return err
		}
		builder.Append(t)
		return nil
	}
}

func timeToTimestamp(value time.Time, _type *arrow.TimestampType) (arrow.Timestamp, error) {
	loc, err := _type.GetZone()
	if err != nil {
		return arrow.Timestamp(0), err
	}
	if loc != nil {
		value = value.In(loc)
	}

	switch _type.Unit {
	case arrow.Second:
		return arrow.Timestamp(value.Unix()), nil
	case arrow.Millisecond:
		return arrow.Timestamp(value.Unix()*1e3 + int64(value.Nanosecond())/1e6), nil
	case arrow.Microsecond:
		return arrow.Timestamp(value.Unix()*1e6 + int64(value.Nanosecond())/1e3), nil
	case arrow.Nanosecond:
		return arrow.Timestamp(value.UnixNano()), nil
	default:
		return arrow.Timestamp(0), fmt.Errorf("unsupported Apache Arrow time unit: %s", _type.Unit.String())
	}
}
