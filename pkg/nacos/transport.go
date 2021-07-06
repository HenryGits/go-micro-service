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
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/sd"
	kitconsul "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"go-micro-service/pkg/nacos/discover"
	"go-micro-service/plugin/pkg/consul"
	"net/http"
	"time"
)

var (
	retryMax = 3
)

type endpoints struct {
	MakeNacosEndpoint endpoint.Endpoint
}

/**
 * 自定义处理器
 */
func MakeHandler(method, path string, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}
	// 创建实例
	// logger 可以使用自己的logger对象 例如zap等
	// instances: 实例对象
	instances := kitconsul.NewInstancer(consul.ConsulClient, logger, "nacosService", []string{"primary"}, true)

	// 获取endpoint可执行对象 与我们直连client.Endpoint()返回的一样
	// 传入 instances: 实例对象  Factory:工厂模式  logger: 日志对象
	endpointer := sd.NewEndpointer(instances, nacosFactory(method, path), logger)

	// 获取所有实例 endpoints
	endpoints, err := endpointer.Endpoints()
	if err == nil {
		_ = level.Info(logger).Log("注册的总实例数", len(endpoints))
	}

	balancer := lb.NewRoundRobin(endpointer)
	retryEndpoint := lb.Retry(retryMax, 5*time.Second, balancer)

	// 路由
	r := mux.NewRouter()
	r.Handle("/nacos/{metaData}", kithttp.NewServer(
		retryEndpoint,
		decodeNacosRequest,
		encodeNacosResponse,
		opts...,
	)).Methods(method)

	return r
}

func decodeNacosRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	name, ok := vars["metaData"]
	if !ok {
		return nil, errors.New("路由异常!")
	}
	return nacos.NacosRequest{
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
	case nacos.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusOK)
	}
}
