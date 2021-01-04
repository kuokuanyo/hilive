package controller

import (
	"hilive/guard"
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/menu"
	"hilive/modules/table"
	"hilive/modules/utils"
	"hilive/template/types"
	"hilive/views/menuviews"
	"html/template"
	"net/http"

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
	ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuNewURL)
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
	formInfo := table.GetMenuFormPanel(h.Conn).GetNewForm(h.Services)

	route := URLRoute{
		PreviousURL: config.Prefix() + h.Config.MenuURL,
		InfoURL:     config.Prefix() + h.Config.MenuNewURL,
		URLPrefix:   config.Prefix(),
		IndexURL:    config.Prefix() + h.Config.IndexURL,
	}

	tmpl, err := template.New("").Funcs(DefaultFuncMap).Parse(menuviews.NewMenuTmpl)
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
