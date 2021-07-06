package nacos

import (
	"github.com/icowan/config"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/util"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	NcClient naming_client.INamingClient
)

type client struct {
	client naming_client.INamingClient
}

type NacosClient interface {
	CreateNacosClient() naming_client.INamingClient
	Subscribe(params vo.SelectAllInstancesParam, ch chan struct{})
	// Register a service with nacos.
	Register(params vo.RegisterInstanceParam) error
	// Deregister a service with nacos.
	Deregister(params vo.DeregisterInstanceParam) error
}

func NewNacosClient(cf *config.Config) (iClient NacosClient, err error) {
	//create clientConfig
	clientConfig := constant.ClientConfig{
		// 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		NamespaceId:         "e786a0bf-2b39-4753-9cbb-314e08f8bb22",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/nacos/log",
		CacheDir:            "/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
		Username:            "nacos",
		Password:            "nacos",
	}

	// At least one ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "localhost",
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
		//{
		//	IpAddr:      "console2.nacos.io",
		//	ContextPath: "/nacos",
		//	Port:        80,
		//	Scheme:      "http",
		//},
	}

	// Another way of create config client for dynamic configuration (recommend)
	//configClient, err := clients.NewConfigClient(
	//	vo.NacosClientParam{
	//		ClientConfig:  &clientConfig,
	//		ServerConfigs: serverConfigs,
	//	},
	//)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	return &client{client: namingClient}, err
}

func (c client) CreateNacosClient() naming_client.INamingClient {
	NcClient = c.client
	return c.client
}

// WatchPrefix implements the etcd Client interface.
// 监听&监控修改，动态配置
func (c *client) Subscribe(params vo.SelectAllInstancesParam, ch chan struct{}) {
	subParam := &vo.SubscribeParam{
		ServiceName: params.ServiceName,
		GroupName:   params.GroupName,
		Clusters:    params.Clusters,
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			println("callback111 return services:%s \n\n", util.ToJsonString(services))

			ch <- struct{}{} // make sure caller invokes GetEntries
			for {
				if err != nil {
					return
				}
				ch <- struct{}{}
			}
		},
	}

	_ = c.client.Subscribe(subParam)
}

func (c client) Register(regParam vo.RegisterInstanceParam) error {
	_, err := c.client.RegisterInstance(regParam)
	return err
}

func (c client) Deregister(degParam vo.DeregisterInstanceParam) error {
	_, err := c.client.DeregisterInstance(degParam)
	return err
}
