/**
* @author: ZHC
* @date: 2021-04-21 12:58:47
* @description: 自定义日志服务
**/
package nacos

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"time"
)

type loggingService struct {
	logger log.Logger
	IService
}

func NewLoggingService(logger log.Logger, s IService) IService {
	return &loggingService{level.Info(logger), s}
}

func (s *loggingService) RegisterServiceInstance(ctx context.Context, metaData string) (rs bool, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"uri", ctx.Value(kithttp.ContextKeyRequestURI),
			"method", "get",
			"name", metaData,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.IService.RegisterServiceInstance(ctx, metaData)
}
