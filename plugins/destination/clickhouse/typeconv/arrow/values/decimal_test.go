package values

import (
	"math"
	"math/rand"
	"testing"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/decimal128"
	"github.com/apache/arrow/go/v12/arrow/decimal256"
	"github.com/shopspring/decimal"
)

func Test_decimal128(t *testing.T) {
	value := decimal.NewFromFloat(rand.Float64()*(math.MaxFloat64/2-1) + rand.Float64())
	for _, _type := range []*arrow.Decimal128Type{
		{Precision: 10, Scale: 0},
		{Precision: 10, Scale: 3},
		{Precision: 10, Scale: 5},
		{Precision: 15, Scale: 0},
		{Precision: 15, Scale: 3},
		{Precision: 15, Scale: 5},
		{Precision: 15, Scale: 10},
		{Precision: 23, Scale: 15},
		{Precision: 32, Scale: 2},
		{Precision: 32, Scale: 15},
		{Precision: 38, Scale: 32},
	} {
		ensureRecord(t, testCase{
			_type:    _type,
			value:    &value,
			expected: decimal128.FromBigInt(value.Coefficient()),
		})
	}
}

func Test_decimal256(t *testing.T) {
	value := decimal.NewFromFloat(rand.Float64()*(math.MaxFloat64/2-1) + rand.Float64())
	for _, _type := range []*arrow.Decimal256Type{
		{Precision: 10, Scale: 0},
		{Precision: 10, Scale: 3},
		{Precision: 10, Scale: 5},
		{Precision: 15, Scale: 0},
		{Precision: 15, Scale: 3},
		{Precision: 15, Scale: 5},
		{Precision: 15, Scale: 10},
		{Precision: 23, Scale: 15},
		{Precision: 32, Scale: 2},
		{Precision: 32, Scale: 15},
		{Precision: 38, Scale: 32},
	} {
		ensureRecord(t, testCase{
			_type:    _type,
			value:    &value,
			expected: decimal256.FromBigInt(value.Coefficient()),
		})
	}
}
