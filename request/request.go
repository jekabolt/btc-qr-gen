package request

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type HTTPClient struct {
	http *resty.Client
}

func NewHTTPClient(c *resty.Client) *HTTPClient {
	hc := &HTTPClient{}
	hc.SetHTTP(c)
	return hc
}

func (c *HTTPClient) HTTP() *resty.Client {
	return c.http
}

func (c *HTTPClient) SetHTTP(hCli *resty.Client) {
	c.http = hCli
}

func (cli *HTTPClient) MakeRequest(method, url string, headers map[string]string, reqBody interface{}) (*resty.Response, error) {
	requestURL := url

	resp := &resty.Response{}
	var err error
	switch method {
	case http.MethodGet:
		resp, err = cli.HTTP().R().EnableTrace().
			SetHeaders(headers).
			Get(requestURL)
		if err != nil {
			return nil, fmt.Errorf("cli.MakeRequest:cli.HTTP:Get:err: [%v]", err.Error())
		}
	case http.MethodDelete:
		resp, err = cli.HTTP().R().EnableTrace().
			SetHeaders(headers).
			Delete(requestURL)
		if err != nil {
			return nil, fmt.Errorf("cli.MakeRequest:HTTP:Delete:err: [%v]", err.Error())
		}
	case http.MethodPost:
		resp, err = cli.HTTP().R().EnableTrace().
			SetHeaders(headers).
			SetBody(reqBody).
			Post(requestURL)
		if err != nil {
			return nil, fmt.Errorf("cli.MakeRequest:HTTP:Post:err: [%v]", err.Error())
		}
	case http.MethodPut:
		resp, err = cli.HTTP().R().EnableTrace().
			SetHeaders(headers).
			SetBody(reqBody).
			Put(requestURL)
		if err != nil {
			return nil, fmt.Errorf("cli.MakeRequest:HTTP:Post:err: [%v]", err.Error())
		}
	}
	return resp, nil
}
