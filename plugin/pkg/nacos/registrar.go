/**
* @author: ZHC
* @date: 2021-04-25 16:43:50
* @description: 注册&注销
**/
package nacos

import (
	"github.com/go-kit/kit/log"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"sync"
	"time"
)

const minHeartBeatTime = 500 * time.Millisecond

// Registrar registers service instance liveness information to etcd.
type Registrar struct {
	client   NacosClient
	regParam vo.RegisterInstanceParam
	degParam vo.DeregisterInstanceParam
	logger   log.Logger

	quitmtx sync.Mutex
	quit    chan struct{}
}

// NewRegistrar returns a etcd Registrar acting on the provided catalog
// registration (service).
func NewRegistrar(client NacosClient, regParam vo.RegisterInstanceParam, degParam vo.DeregisterInstanceParam, logger log.Logger) *Registrar {
	return &Registrar{
		client:   client,
		regParam: regParam,
		degParam: degParam,
		logger:   log.With(logger, "regParam", regParam),
	}
}

// Register implements the sd.Registrar interface. Call it when you want your
// service to be registered in etcd, typically at startup.
func (r *Registrar) Register() {
	if err := r.client.Register(r.regParam); err != nil {
		r.logger.Log("err", err)
	} else {
		r.logger.Log("action", "register")
	}
}

// Deregister implements the sd.Registrar interface. Call it when you want your
// service to be deregistered from etcd, typically just prior to shutdown.
func (r *Registrar) Deregister() {
	if err := r.client.Deregister(r.degParam); err != nil {
		r.logger.Log("err", err)
	} else {
		r.logger.Log("action", "deregister")
	}

	r.quitmtx.Lock()
	defer r.quitmtx.Unlock()
	if r.quit != nil {
		close(r.quit)
		r.quit = nil
	}
}
