package mercari

import (
	"context"
	"net/http"

	"github.com/buyandship/supply-svr/biz/common/trace"
)

func HttpDo(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	ctx, span := trace.StartHTTPOperation(ctx, req)
	defer trace.EndSpan(span, err)

	client := &http.Client{}
	httpRes, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	trace.RecordHTTPResponse(span, httpRes)

	return httpRes, nil
}
