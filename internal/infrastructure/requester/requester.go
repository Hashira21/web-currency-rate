package requester

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Hashira21/currency-rate/internal/models/config"
)

type Requester struct {
	client      http.Client
	method      string
	host        string
	path        string
	queryParams string
}

func New(client *http.Client, provider config.Provider, endpoint string) Requester {
	return Requester{
		client:      *client,
		method:      provider.Endpoints[endpoint].Method,
		host:        provider.Host,
		path:        provider.Endpoints[endpoint].Path,
		queryParams: "",
	}
}

func (req Requester) SetQueryParameters(params url.Values) Requester {
	if len(params) == 0 {
		return req
	}

	encodedUrl := params.Encode()
	req.queryParams = fmt.Sprintf("?%s", encodedUrl)

	return req
}

func (req Requester) DoWithoutBody(ctx context.Context) (*http.Response, error) {
	var reqUrl strings.Builder

	reqUrl.WriteString(req.host)
	reqUrl.WriteString(req.path)
	reqUrl.WriteString(req.queryParams)

	var reader io.Reader

	preparedReq, err := http.NewRequestWithContext(ctx, req.method, reqUrl.String(), reader)
	if err != nil {
		return nil, err
	}

	response, err := req.client.Do(preparedReq)
	if err != nil {
		return nil, err
	}

	return response, err
}
