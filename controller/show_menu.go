package controller

import (
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

// ShowMenu menu GET功能
func (h *Handler) ShowMenu(ctx *gin.Context) {
	// getMenuInfoPanel 取得menu顯示資訊面板
	h.getMenuInfoPanel(ctx, h.Config.MenuURL, "")
}

// getMenuInfoPanel 取得menu前端資訊面板
func (h *Handler) getMenuInfoPanel(ctx *gin.Context, url string, alert string) {
	user := auth.GetUserByMiddleware()
	formInfo := table.GetMenuFormPanel(h.Conn).GetNewForm(h.Services)

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
		FormID   string
		User     models.UserModel
		Menu     *menu.Menu
		Alert    string
		Content  types.FormFields
		Config   *config.Config
		Token    string
		URLRoute URLRoute
	}{
		FormID:   utils.UUID(8),
		User:     user,
		Menu:     menuInfo,
		Alert:    alert,
		Content:  formInfo.FieldList,
		Token:    auth.ConvertInterfaceToTokenService(h.Services.Get("token_csrf_helper")).AddToken(),
		Config:   h.Config,
		URLRoute: route,
	}); err == nil {
		ctx.Status(http.StatusOK)
		ctx.Header("Content-Type", "text/html; charset=utf-8")
	} else {
		ctx.Status(http.StatusOK)
		panic("使用菜單模板發生錯誤")
	}
}
