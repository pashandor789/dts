package log

import (
	"context"
)

type ctxKey int

const (
	fieldsKey ctxKey = 0
)

func Log(ctx context.Context, logger Logger) Logger {
	if f, ok := FieldsFromContext(ctx); ok {
		return logger.WithFields(f)
	}
	return logger
}

func FieldsFromContext(ctx context.Context) (Fields, bool) {
	fields, ok := ctx.Value(fieldsKey).(Fields)
	if ok {
		fieldsCopy := make(map[string]interface{}, len(fields))

		for k, v := range fields {
			fieldsCopy[k] = v
		}

		return fieldsCopy, true
	}

	return Fields{}, false
}
