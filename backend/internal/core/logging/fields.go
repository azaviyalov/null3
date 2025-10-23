package logging

import (
	"log/slog"
	"time"
)

type FieldValue struct {
	v slog.Value
}

func (fv FieldValue) toSlogValue() slog.Value { return fv.v }

type FieldValuer interface {
	ToFieldValue() FieldValue
}

type Field struct {
	Key   string
	Value FieldValue
}

func NewStringValue(s string) FieldValue  { return FieldValue{slog.StringValue(s)} }
func NewUint64Value(u uint64) FieldValue  { return FieldValue{slog.Uint64Value(u)} }
func NewTimeValue(t time.Time) FieldValue { return FieldValue{slog.TimeValue(t)} }

func NewField(key string, val FieldValue) Field { return Field{Key: key, Value: val} }

func CombineFields(fields ...Field) FieldValue {
	slogAttrs := make([]slog.Attr, 0, len(fields))
	for _, f := range fields {
		slogAttrs = append(slogAttrs, slog.Attr{Key: f.Key, Value: f.Value.toSlogValue()})
	}
	return FieldValue{slog.GroupValue(slogAttrs...)}
}
