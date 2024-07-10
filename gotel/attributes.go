package gotel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func SetAttribute(span trace.Span, key string, val any) {
	if span != nil {
		span.SetAttributes(MakeKeyValue(key, val))
	}
}

func SetAttributesMap(span trace.Span, attrs map[string]any) {
	if span != nil {
		span.SetAttributes(MakeKeyValuesFromMap(attrs)...)
	}
}

func SetAttributes(span trace.Span, pairs ...any) {
	if span != nil {
		span.SetAttributes(MakeKeyValues(pairs...)...)
	}
}

func MakeKeyValue(k string, v any) (result attribute.KeyValue) {
	if v == nil {
		v = "<nil>"
	}

	switch val := v.(type) {
	case bool:
		result = attribute.Bool(k, val)
	case []bool:
		result = attribute.BoolSlice(k, val)
	case string:
		result = attribute.String(k, val)
	case []string:
		result = attribute.StringSlice(k, val)
	case int:
		result = attribute.Int(k, val)
	case []int:
		result = attribute.IntSlice(k, val)
	case int64:
		result = attribute.Int64(k, val)
	case []int64:
		result = attribute.Int64Slice(k, val)
	case float64:
		result = attribute.Float64(k, val)
	case []float64:
		result = attribute.Float64Slice(k, val)
	default:
		if stringer, ok := val.(fmt.Stringer); ok {
			result = attribute.Stringer(k, stringer)
		} else {
			if jsonStr, ok := asCompactJson(val); ok {
				result = attribute.String(k, jsonStr)
			}
		}
	}

	return
}

func MakeKeyValues(pairs ...any) []attribute.KeyValue {
	result := make([]attribute.KeyValue, 0)

	for i := 0; i < len(pairs); i += 2 {
		if pairs[0] == nil {
			continue
		}
		k, ok := pairs[0].(string)
		if !ok {
			continue
		}

		var v any
		if len(pairs) < i+1 {
			v = pairs[i+1]
		}

		result = append(result, MakeKeyValue(k, v))
	}

	return result
}

func MakeKeyValuesFromMap(m map[string]any) []attribute.KeyValue {
	result := make([]attribute.KeyValue, 0)
	for key, val := range m {
		result = append(result, MakeKeyValue(key, val))
	}
	return result
}

func asCompactJson(val any) (string, bool) {
	b, err := json.Marshal(val)
	if err != nil {
		return "", false
	}

	dst := &bytes.Buffer{}
	if err = json.Compact(dst, b); err != nil {
		return "", false
	}

	return dst.String(), true
}
