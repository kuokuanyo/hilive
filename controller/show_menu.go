package controller

import (
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/menu"
	"hilive/modules/table"
	"hilive/template/types"
	"hilive/views/menuviews"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// ShowMenu menu GET功能
func (h *Handler) ShowMenu(ctx *gin.Context) {
	// getMenuInfoPanel 取得menu顯示資訊面板
	h.getMenuInfoPanel(ctx, "")
}

// getMenuInfoPanel 取得menu前端資訊面板
func (h *Handler) getMenuInfoPanel(ctx *gin.Context, alert string) {
	user := auth.GetUserByMiddleware()

	menuInfo := menu.GetMenuInformation(user, h.Conn)
	urlPrefix := "/" + h.Config.URLPrefix
	menuEdit := urlPrefix + h.Config.MenuEditURL
	menuDelete := urlPrefix + h.Config.MenuDeleteURL

	tmpl, err := template.New("").Funcs(DefaultFuncMap).Parse(menuviews.MenuTmpl)
	if err != nil {
		panic("使用菜單模板發生錯誤")
	}

	formInfo := table.GetMenuFormPanel(h.Conn).GetNewForm(h.Services)
	formID, err := uuid.NewV4()
	if err != nil {
		panic("產生uuid發生錯誤")
	}

	if err := tmpl.Execute(ctx.Writer, struct {
		FormID       string
		User         models.UserModel
		Menu         *menu.Menu
		Tree         []menu.Item
		Alert        string
		Content      types.FormFields
		NewMenuTitle string
		Title        string
		Description  string
		Previous     string
		Token        string
		MiniLogo     template.HTML
		Logo         template.HTML
		URL          string
		URLPrefix    string
		EditURL      string
		IndexURL     string
		DeleteURL    string
	}{
		FormID:       formID.String(),
		User:         user,
		Menu:         menuInfo,
		Tree:         menuInfo.List,
		Alert:        alert,
		Content:      formInfo.FieldList,
		NewMenuTitle: "新建菜單",
		Title:        "菜單",
		Description:  "菜單處理",
		Previous:     urlPrefix + h.Config.MenuURL,
		Token:        auth.ConvertInterfaceToTokenService(h.Services.Get("token_csrf_helper")).AddToken(),
		MiniLogo:     h.Config.MiniLogo,
		Logo:         h.Config.Logo,
		URL:          urlPrefix + h.Config.MenuURL + "/new",
		URLPrefix:    urlPrefix,
		EditURL:      menuEdit,
		IndexURL:     urlPrefix + h.Config.IndexURL,
		DeleteURL:    menuDelete,
	}); err == nil {
		ctx.Status(http.StatusOK)
		ctx.Header("Content-Type", "text/html; charset=utf-8")
	} else {
		ctx.Status(http.StatusOK)
		panic("使用菜單模板發生錯誤")
	}
}
