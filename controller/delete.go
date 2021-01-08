package controller

import (
	"hilive/guard"
	"hilive/modules/auth"
	"hilive/modules/response"

	"github.com/gin-gonic/gin"
)

// Delete 刪除資料POST功能
func (h *Handler) Delete(ctx *gin.Context) {
	param := guard.GetDeleteParam(ctx)

	err := h.GetTable(ctx, param.Prefix).DeleteData(param.ID)
	if err != nil {
		h.Alert = err.Error()
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.Header("X-PJAX-Url", "/admin/info/"+param.Prefix+param.Param.GetRouteParamStr())
		return
	}

	response.OkWithData(ctx, map[string]interface{}{
		"token": auth.ConvertInterfaceToTokenService(h.Services.Get("token_csrf_helper")).AddToken(),
	})
}
