package customconfiguration

import (
	"net/http"

	"github.com/BachVanHai/my-go-common-lib/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func GetRedisRateLimitConfig(rdb *redis.Client) middleware.RateLimiterConfig {
	return middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		// Sử dụng Identifier là IP hoặc UserID từ JWT
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		Store: &middleware.RateLimiterMemoryStore{
			// Lưu ý: Hiện tại Echo chính thức chưa có RedisStore build-in hoàn hảo
			// Bạn có thể dùng MemoryStore cho từng node hoặc dùng custom store
			// Dưới đây là cách custom ErrorHandler để ép chuẩn Response của bạn
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusTooManyRequests, model.CommonResponseHttp{
				Status:  http.StatusTooManyRequests,
				Message: "Hệ thống đang bận do quá nhiều yêu cầu, vui lòng thử lại sau.",
				Code:    "RATE_LIMIT_REDIS",
			})
		},
	}
}
