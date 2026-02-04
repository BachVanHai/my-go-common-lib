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

			// 1. Mặc định dùng IP
			identifier := c.RealIP()

			// 2. Thử lấy userId từ JWT (Nếu đã qua Middleware JWT trước đó)
			user := c.Get("user").(*jwt.Token) // Echo mặc định lưu token vào key "user"
			if user != nil {
				if claims, ok := user.Claims.(jwt.MapClaims); ok && user.Valid {
					// Lấy userId từ claims (ép kiểu về string để làm key Redis)
					if id, ok := claims["userId"].(string); ok {
						identifier = id
					} else if idFloat, ok := claims["userId"].(float64); ok {
						identifier = strconv.FormatFloat(idFloat, 'f', 0, 64)
					}
				}
			}

			key := "rate_limit:" + identifier

			// 3. Logic Golang thuần với Redis
			// Tăng số đếm
			val, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				return next(c) // Fail-safe: Lỗi Redis thì cho qua
			}

			// Nếu là lần đầu tiên (val == 1), đặt thời gian hết hạn
			if 1 == val {
				rdb.Expire(ctx, key, window)
			}

			// 4. Kiểm tra giới hạn
			if val > limit {
				return c.JSON(http.StatusTooManyRequests, model.CommonResponseHttp{
					Status:  http.StatusTooManyRequests,
					Message: "Tài khoản của bạn đang thực hiện quá nhiều yêu cầu. Vui lòng đợi.",
					Code:    "RATE_LIMIT_USER",
				})
			}

			return next(c)
		}
	}
}
