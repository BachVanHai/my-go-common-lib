package customconfiguration

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/BachVanHai/my-go-common-lib/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func RedisRateLimiter(rdb *redis.Client, limit int64, window time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.Background()

			// 1. Mặc định dùng IP để định danh
			identifier := c.RealIP()

			// 2. Thử lấy userId từ JWT (An toàn, không Panic)
			val := c.Get("user")
			if val != nil {
				if user, ok := val.(*jwt.Token); ok && user.Valid {
					if claims, ok := user.Claims.(jwt.MapClaims); ok {
						// Kiểm tra userId: hỗ trợ cả string và float64 (kiểu mặc định của JSON)
						if idStr, ok := claims["userId"].(string); ok {
							identifier = idStr
						} else if idFloat, ok := claims["userId"].(float64); ok {
							identifier = strconv.FormatFloat(idFloat, 'f', 0, 64)
						}
					}
				}
			}

			key := "rate_limit:" + identifier

			// 3. Logic Redis (Atomic Incr & Expire)
			count, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				// Fail-safe: Nếu Redis có vấn đề, cho qua để không chặn người dùng
				return next(c)
			}

			// Chỉ đặt Expire khi key mới được tạo lần đầu (tránh reset TTL liên tục)
			if count == 1 {
				rdb.Expire(ctx, key, window)
			}

			// 4. Kiểm tra giới hạn
			if count > limit {
				return c.JSON(http.StatusTooManyRequests, model.CommonResponseHttp{
					Status:  http.StatusTooManyRequests,
					Message: "Bạn đã thao tác quá nhanh. Vui lòng thử lại sau.",
					Code:    "RATE_LIMIT_EXCEEDED",
				})
			}

			return next(c)
		}
	}
}
