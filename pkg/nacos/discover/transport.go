/**
* @author: ZHC
* @date: 2021-04-21 12:58:47
* @description: 通信层，可以用各种不同的通信方式，如 HTTP RESTful 接口或者 gRPC 接口（这是个很大的优点，方便切换成任何通信协议）
				数据传输(序列化和反序列化)
**/

/**
 * RPC-client侧
 *	每个接口都会按序安装跟踪、限速、断路器中间件
 *	go-kit的设计不是同一给所有接口安装，而是手动的给每一个接口安装，粒度细了，也多了一点代码量
 */

package nacos

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

const rateBucketNum = 3

type endpoints struct {
	MakeNacosEndpoint endpoint.Endpoint
}

/**
 * 自定义处理器
 */
func DiscoverMakeHandler(svc IService, logger kitlog.Logger) http.Handler {
	//ctx := context.Background()
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	// 终端层
	eps := endpoints{
		MakeNacosEndpoint: MakeNacosEndpoint(svc),
	}
	// 限流中间件
	ems := []endpoint.Middleware{
		NewTokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), rateBucketNum)),
	}

	mw := map[string][]endpoint.Middleware{
		"RegisterServiceInstance": ems,
	}

	for _, m := range mw["RegisterServiceInstance"] {
		eps.MakeNacosEndpoint = m(eps.MakeNacosEndpoint)
	}

	// 路由
	r := mux.NewRouter()
	r.Handle("/nacos/{metaData}", kithttp.NewServer(
		eps.MakeNacosEndpoint,
		decodeNacosRequest,
		encodeNacosResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeNacosRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	name, ok := vars["metaData"]
	if !ok {
		return nil, errors.New("路由异常!")
	}
	return NacosRequest{
		MetaData: name,
	}, nil

}

func encodeNacosResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		return json.NewEncoder(w).Encode(map[string]interface{}{
			"code":  500,
			"error": e,
		})
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

/**
 * 自定义服务错误处理
 */
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusOK)
	}
}
