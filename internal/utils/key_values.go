package utils

import (
	"go.opentelemetry.io/otel/attribute"
	common "go.opentelemetry.io/proto/otlp/common/v1"
)

func KeyValues(attrs []attribute.KeyValue) []*common.KeyValue {
	if len(attrs) == 0 {
		return nil
	}

	out := make([]*common.KeyValue, len(attrs))
	for i, kv := range attrs {
		out[i] = KeyValue(kv)
	}
	return out
}

func KeyValue(kv attribute.KeyValue) *common.KeyValue {
	return &common.KeyValue{Key: string(kv.Key), Value: Value(kv.Value)}
}

func Value(v attribute.Value) *common.AnyValue {
	res := &common.AnyValue{}
	switch v.Type() {
	case attribute.BOOL:
		res.Value = &common.AnyValue_BoolValue{
			BoolValue: v.AsBool(),
		}
	case attribute.BOOLSLICE:
		res.Value = &common.AnyValue_ArrayValue{
			ArrayValue: &common.ArrayValue{
				Values: boolSliceValues(v.AsBoolSlice()),
			},
		}
	case attribute.INT64:
		res.Value = &common.AnyValue_IntValue{
			IntValue: v.AsInt64(),
		}
	case attribute.INT64SLICE:
		res.Value = &common.AnyValue_ArrayValue{
			ArrayValue: &common.ArrayValue{
				Values: int64SliceValues(v.AsInt64Slice()),
			},
		}
	case attribute.FLOAT64:
		res.Value = &common.AnyValue_DoubleValue{
			DoubleValue: v.AsFloat64(),
		}
	case attribute.FLOAT64SLICE:
		res.Value = &common.AnyValue_ArrayValue{
			ArrayValue: &common.ArrayValue{
				Values: float64SliceValues(v.AsFloat64Slice()),
			},
		}
	case attribute.STRING:
		res.Value = &common.AnyValue_StringValue{
			StringValue: v.AsString(),
		}
	case attribute.STRINGSLICE:
		res.Value = &common.AnyValue_ArrayValue{
			ArrayValue: &common.ArrayValue{
				Values: stringSliceValues(v.AsStringSlice()),
			},
		}
	}
	return res
}

func boolSliceValues(vals []bool) []*common.AnyValue {
	converted := make([]*common.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = &common.AnyValue{
			Value: &common.AnyValue_BoolValue{
				BoolValue: v,
			},
		}
	}
	return converted
}

func int64SliceValues(vals []int64) []*common.AnyValue {
	converted := make([]*common.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = &common.AnyValue{
			Value: &common.AnyValue_IntValue{
				IntValue: v,
			},
		}
	}
	return converted
}

func float64SliceValues(vals []float64) []*common.AnyValue {
	converted := make([]*common.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = &common.AnyValue{
			Value: &common.AnyValue_DoubleValue{
				DoubleValue: v,
			},
		}
	}
	return converted
}

func stringSliceValues(vals []string) []*common.AnyValue {
	converted := make([]*common.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = &common.AnyValue{
			Value: &common.AnyValue_StringValue{
				StringValue: v,
			},
		}
	}
	return converted
}
