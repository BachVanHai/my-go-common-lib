package exception

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/BachVanHai/my-go-common-lib/model"
	"github.com/labstack/echo/v4"
)

func CommonErrHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	// Kiểm tra nếu là lỗi của Echo (ví dụ 404, 405)
	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
		message = fmt.Sprintf("%v", he.Message)
	}

	// Bạn có thể thêm các check lỗi logic nghiệp vụ tại đây
	// if myErr, ok := err.(*MyBusinessError); ok { ... }

	// Gửi response theo format chuẩn
	errorResponse := model.CommonResponseHttp{
		Status:  code,
		Message: message,
	}

	if !c.Response().Committed {
		err := c.JSON(code, errorResponse)
		if err != nil {
			return
		}
	}
}
