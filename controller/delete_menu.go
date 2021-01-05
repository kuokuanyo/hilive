package controller

import (
	"hilive/guard"
	"hilive/models"
	"hilive/modules/response"

	"github.com/gin-gonic/gin"
)

// DeleteMenu 刪除菜單POST功能
func (h *Handler) DeleteMenu(ctx *gin.Context) {
	param := guard.GetDeleteMenuParam(ctx)
	// 刪除
	models.SetMenuModelByID(param.ID).SetConn(h.Conn).Delete()
	response.OkWithMsg(ctx, "刪除資料成功")
}
