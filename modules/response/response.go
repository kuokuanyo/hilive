package response

import (
	"hilive/context"
	"net/http"
)

// OkWithData 回傳成功以及data
func OkWithData(ctx *context.Context, data map[string]interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg":  "ok",
		"data": data,
	})
}

// OkWithMsg 回成功以及msg
func OkWithMsg(ctx *context.Context, msg string) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg":  msg,
	})
}

// BadRequest 400錯誤
func BadRequest(ctx *context.Context, msg string) {
	ctx.JSON(http.StatusBadRequest, map[string]interface{}{
		"code": http.StatusBadRequest,
		"msg":  msg,
	})
}

// Error 回傳code:500 and msg
func Error(ctx *context.Context, msg string) {
	ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
		"code": http.StatusInternalServerError,
		"msg":  msg,
	})
}
