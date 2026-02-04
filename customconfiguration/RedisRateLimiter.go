package customconfiguration

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/BachVanHai/my-go-common-lib/model"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

// RedisRateLimiter định nghĩa middleware giới hạn tốc độ dùng Redis
func RedisRateLimiter(rdb *redis.Client, limit int64, window time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.Background()
			// Định danh người dùng bằng IP (hoặc có thể lấy UserID từ JWT)
			key := "rate_limit:" + c.RealIP()

			// Sử dụng lệnh INCR của Redis để đếm
			pipe := rdb.Pipeline()
			incr := pipe.Incr(ctx, key)
			pipe.Expire(ctx, key, window)
			_, err := pipe.Exec(ctx)

			if err != nil {
				return next(c) // Nếu lỗi Redis, cho qua để tránh sập hệ thống
			}

			currentCount := incr.Val()
			if currentCount > limit {
				// ÉP TRẢ VỀ CHUẨN CommonResponseHttp
				return c.JSON(http.StatusTooManyRequests, model.CommonResponseHttp{
					Status:  http.StatusTooManyRequests,
					Message: "Bạn đã vượt quá giới hạn yêu cầu (Max: " + strconv.FormatInt(limit, 10) + "). Thử lại sau.",
					Code:    "TOO_MANY_REQUESTS",
				})
			}

			// Thêm Header để Client biết họ còn bao nhiêu lượt
			c.Response().Header().Set("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
			c.Response().Header().Set("X-RateLimit-Remaining", strconv.FormatInt(limit-currentCount, 10))

			return next(c)
		}
	}
}
