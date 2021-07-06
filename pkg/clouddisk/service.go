/**
* @author: ZHC
* @date: 2021-04-21 12:58:47
* @description: 服务层，实现业务逻辑的地方
**/

package clouddisk

import (
	"context"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
}

func NewService(logger *zap.Logger) IService {
	return &Service{
		logger: logger,
	}
}

type IService interface {
	/**
	 * 上传
	 */
	Upload(ctx context.Context, uploadParam UploadRequest) (result map[string]interface{}, err error)

	/**
	 * 下载
	 */
	Download(ctx context.Context, downloadParam DownloadRequest) (result map[string]interface{}, err error)
}

func (s Service) Upload(ctx context.Context, uploadParam UploadRequest) (result map[string]interface{}, err error) {
	panic("implement me")
}

func (s Service) Download(ctx context.Context, downloadParam DownloadRequest) (result map[string]interface{}, err error) {
	panic("implement me")
}
