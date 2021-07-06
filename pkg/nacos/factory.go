package nacos

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	kithttp "github.com/go-kit/kit/transport/http"
	nacos "go-micro-service/pkg/nacos/discover"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func nacosFactory(method, path string) sd.Factory {
	return func(instance string) (endpoint endpoint.Endpoint, closer io.Closer, err error) {
		if !strings.HasPrefix(instance, "http") {
			instance = "http://" + instance
		}

		tgt, err := url.Parse(instance)
		if err != nil {
			return nil, nil, err
		}
		tgt.Path = path

		var (
			enc kithttp.EncodeRequestFunc
			dec kithttp.DecodeResponseFunc
		)
		enc, dec = encodeNacosRequest, decodeNacosReponse

		return kithttp.NewClient(method, tgt, enc, dec).Endpoint(), nil, nil
	}
}

func encodeNacosRequest(_ context.Context, req *http.Request, itf interface{}) error {
	request := itf.(nacos.NacosRequest)
	path := req.URL.Path
	req.URL.Path = path + "/" + request.MetaData
	println("req.URL:", req.URL.String())
	return nil
}

func decodeNacosReponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response nacos.NacosResponse
	var s map[string]interface{}

	if respCode := resp.StatusCode; respCode >= 400 {
		if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
			return nil, err
		}
		return nil, errors.New(s["error"].(string) + "\n")
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}
