package controller

import (
	"hilive/guard"
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/menu"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"hilive/modules/utils"
	"hilive/template/types"
	"hilive/views/alert"
	"hilive/views/menuviews"
	"html/template"

	"github.com/gin-gonic/gin"
)

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

// ShowEditMenu edit menu GET功能
func (h *Handler) ShowEditMenu(ctx *gin.Context) {
	// 取得middleware驗證後的user
	user := auth.GetUserByMiddleware()
	// GetMenuInformation 透過user取得menu資料表資訊
	menuInfo := menu.GetMenuInformation(user, h.Conn)

	if ctx.Query("id") == "" {
		tmpl, _ := template.New("").Funcs(DefaultFuncMap).Parse(alert.AlertTmpl)
		tmpl.Execute(ctx.Writer, struct {
			User         models.UserModel
			Menu         *menu.Menu
			AlertContent string
			Config       *config.Config
			URLRoute     URLRoute
			IndexURL     string
			URLPrefix    string
		}{
			User:         user,
			Menu:         menuInfo,
			AlertContent: "請填寫id參數",
			Config:       h.Config,
			IndexURL:     config.Prefix() + h.Config.IndexURL,
			URLPrefix:    config.Prefix(),
		})
		return
	}

	formInfo, err := table.GetMenuFormPanel(h.Conn).
		GetDataWithID(parameter.DefaultParameters().SetFieldPKByJoinParam(ctx.Query("id")), h.Services)
	if err != nil {
		tmpl, _ := template.New("").Funcs(DefaultFuncMap).Parse(alert.AlertTmpl)
		tmpl.Execute(ctx.Writer, struct {
			User         models.UserModel
			Menu         *menu.Menu
			AlertContent string
			Config       *config.Config
			URLRoute     URLRoute
			IndexURL     string
			URLPrefix    string
		}{
			User:         user,
			Menu:         menuInfo,
			AlertContent: err.Error(),
			Config:       h.Config,
			IndexURL:     config.Prefix() + h.Config.IndexURL,
			URLPrefix:    config.Prefix(),
		})
	}
	h.showEditMenu(ctx, formInfo, h.Config.MenuEditURL, h.Alert)
}

// showEditMenu /menu/edit模板語法
func (h *Handler) showEditMenu(ctx *gin.Context, formInfo table.FormInfo, url string, alert string) {
	// 取得middleware驗證後的user
	user := auth.GetUserByMiddleware()
	// GetMenuInformation 透過user取得menu資料表資訊
	menuInfo := menu.GetMenuInformation(user, h.Conn).SetActiveClass(url)

	route := URLRoute{
		InfoURL:     config.Prefix() + h.Config.MenuEditURL,
		IndexURL:    config.Prefix() + h.Config.IndexURL,
		URLPrefix:   config.Prefix(),
		PreviousURL: config.Prefix() + h.Config.MenuURL,
	}

	tmpl, err := template.New("").Funcs(DefaultFuncMap).Parse(menuviews.EditMenuTmpl)
	if err != nil {
		panic("使用編輯菜單模板發生錯誤")
	}
	tmpl.Execute(ctx.Writer, struct {
		FormID       string
		User         models.UserModel
		Menu         *menu.Menu
		AlertContent string
		Content      types.FormFields
		Config       *config.Config
		URLRoute     URLRoute
		Token        string
	}{
		FormID:       utils.UUID(8),
		User:         user,
		Menu:         menuInfo,
		AlertContent: alert,
		Content:      formInfo.FieldList,
		Config:       h.Config,
		URLRoute:     route,
		Token:        auth.ConvertInterfaceToTokenService(h.Services.Get("token_csrf_helper")).AddToken(),
	})
	if alert != "" {
		h.Alert = ""
	}
}
