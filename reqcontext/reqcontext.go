package reqcontext

import (
	"context"

	"github.com/google/uuid"
)

type ctxRequestID string

const ctxRequestIDKey = ctxRequestID("requestID")

func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	if requestID == "" {
		requestID = "internal:" + uuid.New().String()
	} else {
		requestID = "user:" + requestID
	}
	return context.WithValue(ctx, ctxRequestIDKey, requestID)
}

func RequestIDFromContext(ctx context.Context) string {
	v := ctx.Value(ctxRequestIDKey)
	if v == nil {
		return ""
	}

	return v.(string)
}
