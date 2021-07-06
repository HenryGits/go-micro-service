/**
* @author: ZHC
* @date: 2021-04-21 12:58:47
* @description: 自定义日志服务
**/
package clouddisk

import (
	"context"
	kithttp "github.com/go-kit/kit/transport/http"
	"go.uber.org/zap"
	"time"
)

type loggingService struct {
	logger *zap.Logger
	IService
}

func NewLoggingService(logger *zap.Logger, s IService) IService {
	return &loggingService{logger, s}
}

func (s *loggingService) RegisterServiceInstance(ctx context.Context, metaData string) (rs bool, err error) {
	defer func(begin time.Time) {
		s.logger.Info("cloudDisk", zap.Any("uri", ctx.Value(kithttp.ContextKeyRequestURI)),
			zap.Any("method", "get"),
			zap.Any("name", metaData),
			zap.Any("took", time.Since(begin)),
			zap.Any("err", err),
		)
	}(time.Now())
	return s.IService.RegisterServiceInstance(ctx, metaData)
}
