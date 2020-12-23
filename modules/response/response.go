package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// OkWithData 回傳成功以及data
func OkWithData(ctx *gin.Context, data map[string]interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg":  "ok",
		"data": data,
	})
}

// OkWithMsg 回成功以及msg
func OkWithMsg(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg":  msg,
	})
}

// BadRequest 400錯誤
func BadRequest(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusBadRequest, map[string]interface{}{
		"code": http.StatusBadRequest,
		"msg":  msg,
	})
}
