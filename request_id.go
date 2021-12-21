package logger

import (
	"context"
	"net/http"
)

const requestIDKey contextKey = "requestID"

func SetRequestID(ctx context.Context, requestID int64) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func getRequestID(ctx context.Context, request *http.Request) int64 {
	if request != nil {
		requestID, ok := ctx.Value(requestIDKey).(int64)
		if ok {
			return requestID
		}
		return 0
	}

	return 0
}
