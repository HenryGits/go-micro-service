module go-micro-service

go 1.16

replace google.golang.org/grpc => google.golang.org/grpc v1.37.0

require (
	github.com/go-kit/kit v0.10.0
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/consul/api v1.8.1
	github.com/icowan/config v0.0.0-20200926110528-b95deb7acc31
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v0.0.0-20190725011945-5c849dd2c51d // indirect
	github.com/nacos-group/nacos-sdk-go v1.0.7
	github.com/satori/go.uuid v1.2.0
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20191205180655-e7c4368fe9dd // indirect
	golang.org/x/net v0.0.0-20210423184538-5f58ad60dda6
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/ini.v1 v1.51.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
