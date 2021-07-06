package nacos

import (
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// Instancer yields instances stored in a certain etcd keyspace. Any kind of
// change in that keyspace is watched and will update the Instancer's Instancers.
type Instancer struct {
	cache       *Cache
	nacosClient NacosClient
	params      vo.SelectAllInstancesParam
	logger      log.Logger
	quitc       chan struct{}
}

// NewInstancer returns an etcd instancer. It will start watching the given
// prefix for changes, and update the subscribers.
//NewInstancer返回一个Nacos instancer。它将开始观察给定的更改的前缀，并更新订阅服务器。
func NewInstancer(c NacosClient, params vo.SelectAllInstancesParam, logger log.Logger) (*Instancer, error) {
	s := &Instancer{
		nacosClient: c,
		params:      params,
		cache:       NewCache(),
		logger:      logger,
		quitc:       make(chan struct{}),
	}

	instances, err := s.nacosClient.CreateNacosClient().SelectAllInstances(params)
	if err == nil {
		logger.Log("params", params, "instances", len(instances))
	} else {
		logger.Log("params", params, "err", err)
	}

	b, err := json.Marshal(instances)
	var its []string
	if json.Unmarshal(b, &its) != nil {
		logger.Log("params", params, "err", "服务注册失败！")
	}
	s.cache.Update(sd.Event{Instances: its, Err: err})

	go s.loop()
	return s, nil
}

func (s *Instancer) loop() {
	ch := make(chan struct{})
	go s.nacosClient.Subscribe(s.params, ch)

	for {
		select {
		case <-ch:
			instances, err := s.nacosClient.CreateNacosClient().SelectAllInstances(s.params)
			if err != nil {
				s.logger.Log("msg", "failed to retrieve entries", "err", err)
				s.cache.Update(sd.Event{Err: err})
				continue
			}
			b, err := json.Marshal(instances)
			var its []string
			if json.Unmarshal(b, &its) != nil {
				s.logger.Log("params", s.params, "err", "服务注册失败！")
			}
			s.cache.Update(sd.Event{Instances: its})

		case <-s.quitc:
			return
		}
	}
}

// Stop terminates the Instancer.
func (s *Instancer) Stop() {
	close(s.quitc)
}

// Register implements Instancer.
func (s *Instancer) Register(ch chan<- sd.Event) {
	s.cache.Register(ch)
}

// Deregister implements Instancer.
func (s *Instancer) Deregister(ch chan<- sd.Event) {
	s.cache.Deregister(ch)
}
