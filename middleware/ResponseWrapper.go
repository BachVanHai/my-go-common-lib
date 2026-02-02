package middleware

import (
	"net/http"

	"github.com/BachVanHai/my-go-common-lib/model"
	"github.com/labstack/echo/v4"
)

func ResponseWrapper(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Cho phép Handler chạy trước
		err := next(c)

		// Nếu đã có response gửi đi rồi (ví dụ trả về file, stream...) thì bỏ qua
		if c.Response().Committed {
			return err
		}

		// Lấy dữ liệu mà Handler đã "gửi" thông qua c.Get hoặc xử lý lỗi
		status := c.Response().Status
		var data any
		message := "Success"
		code := ""

		// Nếu Handler trả về lỗi
		if err != nil {
			if he, ok := err.(*echo.HTTPError); ok {
				status = he.Code
				message = he.Message.(string)
			} else {
				status = http.StatusInternalServerError
				message = err.Error()
			}
			code = "INTERNAL_ERROR" // Hoặc logic lấy mã lỗi riêng
		} else {
			// Lấy data từ context nếu Handler có set
			data = c.Get("result")
			if msg := c.Get("message"); msg != nil {
				message = msg.(string)
			}
		}

		// Đóng gói vào struct chuẩn của bạn
		response := model.CommonResponseHttp{
			Status:          status,
			Message:         message,
			Code:            code,
			Data:            data,
			ClientMessageId: c.Request().Header.Get("clientMessageId"),
		}

		// Gửi response chuẩn JSON
		return c.JSON(status, response)
	}
}
