/**
* @author: ZHC
* @date: 2021-04-21 12:58:47
* @description: 中间件
**/

package clouddisk

import (
	"context"
	"errors"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
)

var ErrLimitExceed = errors.New("限流~ Rate limit exceed!")

/**
 * 限流器（令牌桶）
 */
func NewTokenBucketLimitter(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

/**
 * 熔断器
 */
func Hystrix(commandName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			var resp interface{}
			if resp = circuitbreaker.Hystrix(commandName); resp != nil {
				return resp, err
			}
			return next(ctx, request)
		}
	}
}
