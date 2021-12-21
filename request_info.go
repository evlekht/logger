package logger

import (
	"context"
	"net/http"
)

type contextKey string

const requestInfoKey contextKey = "request_info"

type requestInfo struct {
	RequestID int64
	Headers   map[string][]string
	Method    string
	URL       string
	Body      []byte
}

func MakeRequestInfoContext(ctx context.Context, request *http.Request) context.Context {
	reqInfo := &requestInfo{}
	reqInfo.RequestID = getRequestID(ctx, request)

	if request != nil {
		reqInfo.URL = request.URL.String()
		reqInfo.Method = request.Method
		reqInfo.Headers = request.Header

		if request.GetBody != nil {
			if body, err := request.GetBody(); err != nil {
				reqInfo.Body = []byte("error getting body")
			} else {
				if _, err := body.Read(reqInfo.Body); err != nil {
					reqInfo.Body = []byte("error reading body")
				}
			}
		}
	}

	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, requestInfoKey, reqInfo)
}

func getRequestInfo(ctx context.Context) *requestInfo {
	if reqInfo, ok := ctx.Value(requestInfoKey).(*requestInfo); ok {
		return reqInfo
	}
	return nil
}
