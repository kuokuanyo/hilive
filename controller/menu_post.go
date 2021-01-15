package controller

import (
	"hilive/context"
	"hilive/guard"
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/menu"
	"hilive/modules/parameter"
	"hilive/modules/response"
	"hilive/modules/table"
)

// NewMenu 新建菜單POST功能
func (h *Handler) NewMenu(ctx *context.Context) {
	// 取得Parameters["new_menu_param"]
	param := guard.GetNewMenuParam(ctx)
	if param.Alert != "" {
		h.getMenuInfoPanel(ctx, param.Alert)
		ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
		ctx.AddHeader("X-PJAX-Url", config.Prefix()+h.config.MenuNewURL)
		return
	}

	// 取得目前登入用戶(Context.UserValue["user"])並轉換成UserModel
	user := auth.Auth(ctx)

	// 新建菜單
	menuModel, err := models.DefaultMenuModel().SetConn(h.conn).
		New(param.Title, param.Icon, param.URL, param.Header, param.ParentID, (menu.GetMenuInformation(user, h.conn)).MaxOrder+1)
	if err != nil {
		h.showNewMenu(ctx, "新建菜單發生錯誤，請重新操作")
		ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
		ctx.AddHeader("X-PJAX-Url", config.Prefix()+h.config.MenuNewURL)
		return
	}

	for _, roleID := range param.Roles {
		_, err = menuModel.AddRole(roleID)
		if err != nil {
			h.showNewMenu(ctx, "新建菜單角色發生錯誤，請重新操作")
			ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
			ctx.AddHeader("X-PJAX-Url", config.Prefix()+h.config.MenuNewURL)
			return
		}
	}
	// 增加MaxOrder
	menu.GetMenuInformation(user, h.conn).MaxOrder++

	h.getMenuInfoPanel(ctx, "")
	ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
	ctx.AddHeader("X-PJAX-Url", config.Prefix()+h.config.MenuURL)
}

// EditMenu 編輯菜單POST功能
func (h *Handler) EditMenu(ctx *context.Context) {
	param := guard.GetEditMenuParam(ctx)
	if param.Alert != "" {
		h.getMenuInfoPanel(ctx, param.Alert)
		ctx.AddHeader("X-PJAX-Url", config.Prefix()+h.config.MenuURL)
		return
	}

	// 建立MenuModel
	menuModel := models.SetMenuModelByID(param.ID).SetConn(h.conn)

	// 先刪除所有角色
	err := menuModel.DeleteRoles()
	if err != nil {
		if err.Error() != "沒有影響任何資料" {
			formInfo, _ := table.GetMenuPanel(h.conn).
				GetDataWithID(parameter.DefaultParameters().SetPKs(ctx.Query("id")), h.services)
			h.showEditMenu(ctx, formInfo, "刪除角色發生錯誤，請重新操作")
			ctx.AddHeader("X-PJAX-Url", config.Prefix()+h.config.MenuURL)
			return
		}
	}

	// 新建角色
	for _, roleID := range param.Roles {
		_, err = menuModel.AddRole(roleID)
		if err != nil {
			formInfo, _ := table.GetMenuPanel(h.conn).
				GetDataWithID(parameter.DefaultParameters().SetPKs(ctx.Query("id")), h.services)
			h.showEditMenu(ctx, formInfo, "新增角色發生錯誤，請重新操作")
			ctx.AddHeader("X-PJAX-Url", config.Prefix()+h.config.MenuURL)
			return
		}
	}

	// 更新資料
	_, err = menuModel.Update(param.Title, param.Icon, param.URL, param.Header, param.ParentID)
	if err != nil {
		if err.Error() != "沒有影響任何資料" {
			formInfo, _ := table.GetMenuPanel(h.conn).
				GetDataWithID(parameter.DefaultParameters().SetPKs(ctx.Query("id")), h.services)
			h.showEditMenu(ctx, formInfo, "更新角色發生錯誤，請重新操作")
			ctx.AddHeader("X-PJAX-Url", config.Prefix()+h.config.MenuURL)
			return
		}
	}

	h.getMenuInfoPanel(ctx, "")
	ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
	ctx.AddHeader("X-PJAX-Url", config.Prefix()+h.config.MenuURL)
}

// DeleteMenu 刪除菜單POST功能
func (h *Handler) DeleteMenu(ctx *context.Context) {
	param := guard.GetDeleteMenuParam(ctx)
	// 刪除
	models.SetMenuModelByID(param.ID).SetConn(h.conn).Delete()
	response.OkWithMsg(ctx, "刪除資料成功")
}
