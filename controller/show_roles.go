package controller

import (
	"hilive/modules/parameter"
	"hilive/modules/table"

	"github.com/gin-gonic/gin"
)

// ShowRolesInfo 前端資訊頁面
func (h *Handler) ShowRolesInfo(ctx *gin.Context) {
	// 取得角色面板資訊
	panel := table.GetRolesInfoPanel(h.Conn)
	// 設置頁面資訊
	params := parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize)
	// 取得頁面資料後並執行前端模板語法
	h.showTable(ctx, params, panel, h.Config.RolesURL, "roles")
}
