package logger

import (
	"context"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
)

const requestInfoKey = "request_info" // no custom type cause of fasthttp.RequestCtx

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

	switch requestContext := ctx.(type) {
	case *fasthttp.RequestCtx:
		reqInfo.Body = requestContext.Request.Body()

		headers := make(map[string][]string)
		requestContext.Request.Header.VisitAll(func(key []byte, value []byte) {
			headers[string(key)] = strings.Split(string(value), ",")
		})

		reqInfo.URL = string(requestContext.URI().Path())
		reqInfo.Method = string(requestContext.Method())
		reqInfo.Headers = headers

		requestContext.SetUserValue(requestInfoKey, reqInfo)
		return requestContext

	default:
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
}

func getRequestInfo(ctx context.Context) *requestInfo {
	if reqInfo, ok := ctx.Value(requestInfoKey).(*requestInfo); ok {
		return reqInfo
	}
	return nil
}
