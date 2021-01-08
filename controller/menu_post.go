package controller

import (
	"hilive/guard"
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/menu"
	"hilive/modules/response"

	"github.com/gin-gonic/gin"
)

// NewMenu 新建菜單POST功能
func (h *Handler) NewMenu(ctx *gin.Context) {
	param := guard.GetNewMenuParam(ctx)
	if param.Alert != "" {
		h.Alert = param.Alert
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuURL)
		return
	}
	// GetUserByMiddleware 取得middleware驗證後的user
	user := auth.GetUserByMiddleware()

	// 新建菜單
	menuModel, err := models.DefaultMenuModel().SetConn(h.Conn).
		New(param.Title, param.Icon, param.URL, param.Header, param.ParentID, (menu.GetMenuInformation(user, h.Conn)).MaxOrder+1)
	if err != nil {
		h.Alert = "新建菜單發生錯誤"
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuNewURL)
		return
	}

	for _, roleID := range param.Roles {
		_, err = menuModel.AddRole(roleID)
		if err != nil {
			h.Alert = "新建角色發生錯誤"
			ctx.Header("Content-Type", "text/html; charset=utf-8")
			ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuNewURL)
			return
		}
	}
	// 增加MaxOrder
	menu.GetMenuInformation(user, h.Conn).MaxOrder++
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuURL)
}

// EditMenu 編輯菜單POST功能
func (h *Handler) EditMenu(ctx *gin.Context) {
	param := guard.GetEditMenuParam(ctx)
	if param.Alert != "" {
		h.Alert = param.Alert
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuURL)
		return
	}

	// 建立MenuModel
	menuModel := models.SetMenuModelByID(param.ID).SetConn(h.Conn)

	// 先刪除所有角色
	err := menuModel.DeleteRoles()
	if err != nil {
		if err.Error() != "沒有影響任何資料" {
			h.Alert = "刪除角色發生錯誤"
			ctx.Header("Content-Type", "text/html; charset=utf-8")
			ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuURL)
			return
		}
	}

	// 新建角色
	for _, roleID := range param.Roles {
		_, err = menuModel.AddRole(roleID)
		if err != nil {
			h.Alert = "新建角色發生錯誤"
			ctx.Header("Content-Type", "text/html; charset=utf-8")
			ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuURL)
			return
		}
	}

	// 更新資料
	_, err = menuModel.Update(param.Title, param.Icon, param.URL, param.Header, param.ParentID)
	if err != nil {
		if err.Error() != "沒有影響任何資料" {
			h.Alert = "更新菜單資料發生錯誤"
			ctx.Header("Content-Type", "text/html; charset=utf-8")
			ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuURL)
			return
		}
	}

	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuURL)
}

// DeleteMenu 刪除菜單POST功能
func (h *Handler) DeleteMenu(ctx *gin.Context) {
	param := guard.GetDeleteMenuParam(ctx)
	// 刪除
	models.SetMenuModelByID(param.ID).SetConn(h.Conn).Delete()
	response.OkWithMsg(ctx, "刪除資料成功")
}
