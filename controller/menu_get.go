package controller

import (
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/menu"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"hilive/modules/utils"
	"hilive/template/types"
	"hilive/views/alert"
	"hilive/views/form"
	"hilive/views/menuviews"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ShowMenu menu GET功能
func (h *Handler) ShowMenu(ctx *gin.Context) {
	// getMenuInfoPanel 取得menu顯示資訊面板
	h.getMenuInfoPanel(ctx, h.Config.MenuURL, h.Alert)
}

// getMenuInfoPanel 取得menu前端資訊面板
func (h *Handler) getMenuInfoPanel(ctx *gin.Context, url string, alert string) {
	user := auth.GetUserByMiddleware()
	formInfo := table.GetMenuPanel(h.Conn).GetNewForm(h.Services)

	menuInfo := menu.GetMenuInformation(user, h.Conn).SetActiveClass(url)

	route := URLRoute{
		PreviousURL: config.Prefix() + h.Config.MenuURL,
		URLPrefix:   config.Prefix(),
		InfoURL:     config.Prefix() + h.Config.MenuURL + "/new",
		IndexURL:    config.Prefix() + h.Config.IndexURL,
		EditURL:     config.Prefix() + h.Config.MenuEditURL,
		DeleteURL:   config.Prefix() + h.Config.MenuDeleteURL,
	}

	tmpl, err := template.New("").Funcs(DefaultFuncMap).Parse(menuviews.MenuTmpl)
	if err != nil {
		panic("使用菜單模板發生錯誤")
	}
	if err := tmpl.Execute(ctx.Writer, struct {
		FormID       string
		User         models.UserModel
		Menu         *menu.Menu
		AlertContent string
		Content      types.FormFields
		Config       *config.Config
		Token        string
		URLRoute     URLRoute
	}{
		FormID:       utils.UUID(8),
		User:         user,
		Menu:         menuInfo,
		AlertContent: alert,
		Content:      formInfo.FieldList,
		Token:        auth.ConvertInterfaceToTokenService(h.Services.Get("token_csrf_helper")).AddToken(),
		Config:       h.Config,
		URLRoute:     route,
	}); err == nil {
		ctx.Status(http.StatusOK)
		ctx.Header("Content-Type", "text/html; charset=utf-8")
	} else {
		panic("使用新建菜單模板發生錯誤")
	}
	if alert != "" {
		h.Alert = ""
	}
}

// ShowNewMenu new menu GET功能
func (h *Handler) ShowNewMenu(ctx *gin.Context) {
	h.showNewMenu(ctx, h.Config.MenuNewURL, h.Alert)
}

// showNewMenu new menu面板
func (h *Handler) showNewMenu(ctx *gin.Context, url string, alert string) {
	// GetUserByMiddleware 取得middleware驗證後的user
	user := auth.GetUserByMiddleware()
	// GetMenuInformation 透過user取得menu資料表資訊
	menuInfo := menu.GetMenuInformation(user, h.Conn).SetActiveClass(url)
	// 取得表單資訊
	formInfo := table.GetMenuPanel(h.Conn).GetNewForm(h.Services)

	route := URLRoute{
		PreviousURL: config.Prefix() + h.Config.MenuURL,
		InfoURL:     config.Prefix() + h.Config.MenuNewURL,
		URLPrefix:   config.Prefix(),
		IndexURL:    config.Prefix() + h.Config.IndexURL,
	}

	tmpl, err := template.New("").Funcs(DefaultFuncMap).Parse(form.FormTmpl)
	if err != nil {
		panic("使用新建菜單模板發生錯誤")
	}
	if err := tmpl.Execute(ctx.Writer, struct {
		FormID       string
		User         models.UserModel
		Content      types.FormFields
		AlertContent string
		Menu         *menu.Menu
		Config       *config.Config
		URLRoute     URLRoute
		Token        string
	}{
		FormID:       utils.UUID(8),
		User:         user,
		Content:      formInfo.FieldList,
		AlertContent: alert,
		Menu:         menuInfo,
		Config:       h.Config,
		URLRoute:     route,
		Token:        auth.ConvertInterfaceToTokenService(h.Services.Get("token_csrf_helper")).AddToken(),
	}); err == nil {
		ctx.Status(http.StatusOK)
		ctx.Header("Content-Type", "text/html; charset=utf-8")
	} else {
		panic("使用新建菜單模板發生錯誤")
	}
	if alert != "" {
		h.Alert = ""
	}
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

	formInfo, err := table.GetMenuPanel(h.Conn).
		GetDataWithID(parameter.DefaultParameters().SetPKs(ctx.Query("id")), h.Services)
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

// showEditMenu edit menu 模板語法
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

	tmpl, err := template.New("").Funcs(DefaultFuncMap).Parse(form.FormTmpl)
	if err != nil {
		panic("使用編輯菜單模板發生錯誤")
	}
	tmpl.Execute(ctx.Writer, struct {
		FormID       string
		User         models.UserModel
		Menu         *menu.Menu
		AlertContent string
		FormInfo     table.FormInfo
		Config       *config.Config
		URLRoute     URLRoute
		Token        string
	}{
		FormID:       utils.UUID(8),
		User:         user,
		Menu:         menuInfo,
		AlertContent: alert,
		FormInfo:     formInfo,
		Config:       h.Config,
		URLRoute:     route,
		Token:        auth.ConvertInterfaceToTokenService(h.Services.Get("token_csrf_helper")).AddToken(),
	})
	if alert != "" {
		h.Alert = ""
	}
}
