/**
* @author: ZHC
* @date: 2021-04-21 12:58:47
* @description: 终端层，主要实现各种接口的 handler，负责 req／resp 格式的转换
**/

package nacos

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

/**
 * 请求参数
 */
type NacosRequest struct {
	MetaData string `json:"metadata"`
}

type NacosResponse struct {
	//Success bool   `json:"success"`
	Code  int                    `json:"code"`
	Token string                 `json:"token"`
	Data  map[string]interface{} `json:"data"`
	Err   error                  `json:"error"`
}

func (r NacosResponse) error() error { return r.Err }

func MakeNacosEndpoint(s IService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//req := request.(NacosRequest)
		var code int = 200
		//rs, err := s.RegisterServiceInstance(ctx, req.MetaData)
		rs, err := s.GetTodo(ctx)
		if err != nil {
			code = 500
		}
		result := map[string]interface{}{}
		result["data"] = rs
		return NacosResponse{Code: code, Err: err, Data: result}, err
	}
}
