package controller

import (
	"hilive/guard"

	"github.com/gin-gonic/gin"
)

// EditForm 更新資料POST功能
func (h *Handler) EditForm(ctx *gin.Context) {
	param := guard.GetEditForm(ctx)

	err := param.Panel.UpdateData(param.MultiForm.Value)
	if err != nil {
		h.Alert = err.Error()
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.Header("X-PJAX-Url", param.Path+"/edit?__edit_pk="+param.ID)
		return
	}

	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Header("X-PJAX-Url", param.Path+param.Param.DeletePK().GetRouteParamStr())
}
