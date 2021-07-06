package server

import (
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitConsul "github.com/go-kit/kit/sd/consul"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	nacos2 "go-micro-service/pkg/nacos/discover"
	"go-micro-service/plugin/pkg/consul"
	"go-micro-service/plugin/pkg/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const (
	DefaultHttpPort      = ":8880"
	DefaultConfigPath    = "app.cfg"
	DefaultKubeConfig    = "config.yaml"
	DefaultAdminEmail    = "kplcloud@nsini.com"
	DefaultAdminPassword = "admin"
	DefaultInitDBSQL     = "./database/kplcloud.sql"
)

var (
	log kitlog.Logger
)

/**
 *	这是演示 grpc-web+vue+Nginx 搭建一个简单todo示例
 */
//func RunTodo() {
//	lis, err := net.Listen("tcp", DefaultHttpPort)
//	if err != nil {
//		logger.Log("failed to listen: %v", err)
//	}
//
//	s := todo.Server{}
//	grpcServer := grpc.NewServer()
//	// attach the todo service to the server
//	todo.RegisterTodoServiceServer(grpcServer, &s)
//
//	// Register reflection service on gRPC server.
//	reflection.Register(grpcServer)
//
//	if err := grpcServer.Serve(lis); err != nil {
//		logger.Log("failed to serve: %s", err)
//	} else {
//		logger.Log("Server started successfully")
//	}
//
//}
//
//func nacosRun() {
//	//nacosService := nacos.NewService()
//	//
//	//mux := http.NewServeMux()
//	//mux.Handle("/namespace", nacos.MakeHandler(nacosService, httpLogger))
//
//	nc, err := nacos.NewNacosClient(nil)
//	if err != nil {
//		fmt.Println(err)
//	}
//	nacosClient := nc.CreateNacosClient()
//	success, err := nacosClient.RegisterInstance(vo.RegisterInstanceParam{
//		Ip:          "127.0.0.1",
//		Port:        8848,
//		ServiceName: "todo.todoService",
//		Weight:      10,
//		Enable:      true,
//		Healthy:     true,
//		Ephemeral:   true,
//		Metadata:    map[string]string{"idc": "shanghai"},
//		ClusterName: "cluster-a", // 默认值DEFAULT
//		GroupName:   "group-a",   // 默认值DEFAULT_GROUP
//	})
//
//	if err != nil {
//		fmt.Println("RegisterInstance异常：", err)
//	}
//
//	// 获取服务信息：GetService
//	services, err := nacosClient.GetService(vo.GetServiceParam{
//		ServiceName: "todo.todoService",
//		Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
//		GroupName:   "group-a",             // 默认值DEFAULT_GROUP
//	})
//
//	fmt.Println("success: ", success)
//	fmt.Println("services: ", services)
//
//	// SelectOneHealthyInstance将会按加权随机轮询的负载均衡策略返回一个健康的实例
//	// 实例必须满足的条件：health=true,enable=true and weight>0
//	instance, err := nacosClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
//		ServiceName: "todo.todoService",
//		GroupName:   "group-a",             // 默认值DEFAULT_GROUP
//		Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
//	})
//
//	fmt.Println("instance: ", instance.InstanceId)
//
//	//Subscribe key=serviceName+groupName+cluster
//	//Note:We call add multiple SubscribeCallback with the same key.
//	param := &vo.SubscribeParam{
//		ServiceName: "todo.todoService",
//		GroupName:   "group-a",
//		Clusters:    []string{"cluster-a"},
//		SubscribeCallback: func(services []model.SubscribeService, err error) {
//			fmt.Printf("callback111 return services:%s \n\n", util.ToJsonString(services))
//		},
//	}
//	_ = nacosClient.Subscribe(param)
//
//	RunTodo()
//}

func Run() {
	logger := logging.SetLogging(zap.NewExample(), "", zapcore.InfoLevel)

	svcHost := "localhost"
	svcPort := ":8500"
	_, err := consul.NewConsulClient("http://" + svcHost + svcPort)
	if err != nil {
		fmt.Println(err)
	}

	nacosSvc := nacos2.NewService(logger)
	nacosSvc = nacos2.NewLoggingService(logger, nacosSvc)

	mux := http.NewServeMux()
	mux.Handle("/nacos/", nacos2.DiscoverMakeHandler(nacosSvc, logger))

	handlers := make(map[string]string, 3)
	http.Handle("/health", health())

	http.Handle("/", accessControl(mux, logger, handlers))

	// 设置Consul对服务健康检查的参数
	check := api.AgentServiceCheck{
		HTTP:     "http://" + svcHost + DefaultHttpPort + "/health",
		Interval: "10s",
		Timeout:  "1s",
		Notes:    "Consul check service health status.",
	}

	port, _ := strconv.Atoi(DefaultHttpPort[1:])

	//设置微服务想Consul的注册信息
	reg := api.AgentServiceRegistration{
		ID:      "nacosService" + uuid.New().String(),
		Name:    "nacosService",
		Address: "localhost",
		Port:    port,
		Tags:    []string{"primary"},
		Check:   &check,
	}

	//client, err :=api.NewClient()
	//client.Agent().ServiceRegister(&reg)

	errs := make(chan error, 2)
	c := make(chan os.Signal)
	go func() {
		_ = level.Info(logger).Log("transport", "http")
		//创建注册
		registrar := kitConsul.NewRegistrar(consul.ConsulClient, &reg, logger)
		//启动注册服务
		registrar.Register()
		//errs <- http.ListenAndServe(httpAddr, addCors())
		errs <- http.ListenAndServe(DefaultHttpPort, nil)
	}()
	go func() {
		// 接收到信号
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	_ = level.Error(logger).Log("terminated", <-errs)

}

func accessControl(h http.Handler, logger kitlog.Logger, headers map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, val := range headers {
			w.Header().Set(key, val)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Connection", "keep-alive")

		if r.Method == "OPTIONS" {
			return
		}

		//requestId := r.Header.Get("X-Request-Id")
		_ = level.Info(logger).Log("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)
		h.ServeHTTP(w, r)
	})

}

func health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
}
