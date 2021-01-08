package controller

import (
	"hilive/guard"

	"github.com/gin-gonic/gin"
)

// NewForm 新增資料POST功能
func (h *Handler) NewForm(ctx *gin.Context) {
	param := guard.GetNewForm(ctx)

	err := param.Panel.InsertData(param.MultiForm.Value)
	if err != nil {
		h.Alert = err.Error()
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.Header("X-PJAX-Url", param.Path+"/new")
		return
	}

	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Header("X-PJAX-Url", param.Path+param.Param.GetRouteParamStr())
}
