package controller

import (
	"bytes"
	"hilive/context"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"net/http"
)

// ShowMenu menu GET功能
func (h *Handler) ShowMenu(ctx *context.Context) {
	// getMenuInfoPanel 取得menu顯示資訊面板
	h.getMenuInfoPanel(ctx, "")
}

// getMenuInfoPanel 取得menu前端資訊面板
func (h *Handler) getMenuInfoPanel(ctx *context.Context, alert string) {
	// Auth 取得目前登入用戶(Context.UserValue["user"])並轉換成UserModel
	user := auth.Auth(ctx)

	// GetNewForm 處理並設置表單欄位細節資訊(允許增加的表單欄位)
	formInfo := table.GetMenuPanel(h.conn).GetNewForm(h.services)

	route := URLRoute{
		PreviousURL: config.Prefix() + h.config.MenuURL,
		URLPrefix:   config.Prefix(),
		InfoURL:     config.Prefix() + h.config.MenuURL + "/new",
		IndexURL:    config.Prefix() + h.config.IndexURL,
		EditURL:     config.Prefix() + h.config.MenuEditURL,
		DeleteURL:   config.Prefix() + h.config.MenuDeleteURL,
	}

	formInfo.HideBackButton = true
	buf := h.ExecuteForm(ctx, user, formInfo, route, alert, "menu")
	ctx.HTML(http.StatusOK, buf.String())
}

// ShowNewMenu new menu GET功能
func (h *Handler) ShowNewMenu(ctx *context.Context) {
	h.showNewMenu(ctx, "")
}

// showNewMenu new menu面板
func (h *Handler) showNewMenu(ctx *context.Context, alert string) {
	// Auth 取得目前登入用戶(Context.UserValue["user"])並轉換成UserModel
	user := auth.Auth(ctx)

	// 取得表單資訊
	formInfo := table.GetMenuPanel(h.conn).GetNewForm(h.services)

	route := URLRoute{
		PreviousURL: config.Prefix() + h.config.MenuURL,
		InfoURL:     config.Prefix() + h.config.MenuNewURL,
		URLPrefix:   config.Prefix(),
		IndexURL:    config.Prefix() + h.config.IndexURL,
	}

	buf := h.ExecuteForm(ctx, user, formInfo, route, alert, "form")
	ctx.HTML(http.StatusOK, buf.String())
}

// ShowEditMenu edit menu GET功能
func (h *Handler) ShowEditMenu(ctx *context.Context) {
	// Auth 取得目前登入用戶(Context.UserValue["user"])並轉換成UserModel
	user := auth.Auth(ctx)

	buf := new(bytes.Buffer)
	if ctx.Query("id") == "" {
		buf = h.ExecuteAlert(ctx, user, "發生錯誤:請填寫需要編輯的id參數")
		ctx.HTML(http.StatusOK, buf.String())
		return
	}

	formInfo, err := table.GetMenuPanel(h.conn).
		GetDataWithID(parameter.DefaultParameters().SetPKs(ctx.Query("id")), h.services)
	if err != nil {
		buf = h.ExecuteAlert(ctx, user, err.Error())
		ctx.HTML(http.StatusOK, buf.String())
		return
	}
	h.showEditMenu(ctx, formInfo, "")
}

// showEditMenu edit menu 模板語法
func (h *Handler) showEditMenu(ctx *context.Context, formInfo table.FormInfo, alert string) {
	// Auth 取得目前登入用戶(Context.UserValue["user"])並轉換成UserModel
	user := auth.Auth(ctx)

	route := URLRoute{
		InfoURL:     config.Prefix() + h.config.MenuEditURL,
		IndexURL:    config.Prefix() + h.config.IndexURL,
		URLPrefix:   config.Prefix(),
		PreviousURL: config.Prefix() + h.config.MenuURL,
	}

	buf := h.ExecuteForm(ctx, user, formInfo, route, alert, "form")
	ctx.HTML(http.StatusOK, buf.String())
	return
}
