/**
* @author: ZHC
* @date: 2021-04-21 12:58:47
* @description: 服务层，实现业务逻辑的地方
**/

package nacos

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go-micro-service/plugin/pkg/nacos"
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
)

type IService interface {
	/**
	 * 注册Nacos服务
	 */
	RegisterServiceInstance(ctx context.Context, metaData string) (rs bool, err error)
	/**
	 * 获取服务信息
	 */
	GetService(ctx context.Context, param vo.GetServiceParam) (rs string, err error)

	GetTodo(ctx context.Context) (rs string, err error)
}

type Service struct {
	logger log.Logger
}

func NewService(logger log.Logger) IService {
	return &Service{
		logger: logger,
	}
}

/**
 * 注册Nacos服务
 */
func (c *Service) RegisterServiceInstance(ctx context.Context, metaData string) (rs bool, err error) {
	client := nacos.NcClient
	param := vo.RegisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        8848,
		ServiceName: "todo.todoService",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"metadata": metaData},
		ClusterName: "cluster-a", // 默认值DEFAULT
		GroupName:   "group-a",   // 默认值DEFAULT_GROUP
	}
	return client.RegisterInstance(param)
}

/**
 * 获取服务信息
 */
func (c *Service) GetService(ctx context.Context, param vo.GetServiceParam) (rs string, err error) {
	services, err := nacos.NcClient.GetService(param)
	jsonServices, err := json.Marshal(services)
	return string(jsonServices), err
}

func (c *Service) GetTodo(ctx context.Context) (rs string, err error) {
	return "ok!", nil
}
