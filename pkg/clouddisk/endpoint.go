/**
* @author: ZHC
* @date: 2021-04-21 12:58:47
* @description: 网盘终端层，主要实现各种接口的 handler，负责 req／resp 格式的转换
**/

package clouddisk

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

/**
 * 请求参数
 */
type UploadRequest struct {
	MetaData string `json:"metadata"`
}

type CloudDiskResponse struct {
	//Success bool   `json:"success"`
	Code  int                    `json:"code"`
	Token string                 `json:"token"`
	Data  map[string]interface{} `json:"data"`
	Err   error                  `json:"error"`
}

type DownloadRequest struct {
	MetaData string `json:"metadata"`
}

func (r UploadResponse) error() error { return r.Err }

func MakeuploadEndpoint(s IService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//req := request.(uploadRequest)
		var code int = 200
		//rs, err := s.RegisterServiceInstance(ctx, req.MetaData)
		rs, err := s.GetTodo(ctx)
		if err != nil {
			code = 500
		}
		result := map[string]interface{}{}
		result["data"] = rs
		return UploadResponse{Code: code, Err: err, Data: result}, err
	}
}
