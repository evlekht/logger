package logger

import (
	"context"
	"net/http"

	"github.com/valyala/fasthttp"
)

const requestIDKey = "requestID" // no custom type cause of fasthttp.RequestCtx

func SetRequestID(ctx context.Context, requestID int64) context.Context {
	switch requestContext := ctx.(type) {
	case *fasthttp.RequestCtx:
		requestContext.SetUserValue(requestIDKey, requestID)
		return requestContext
	default:
		return context.WithValue(ctx, requestIDKey, requestID)
	}
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
