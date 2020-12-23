package controller

import (
	"hilive/guard"
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

// NewMenu 新建菜單POST功能
func (h *Handler) NewMenu(ctx *gin.Context) {
	param := guard.GetNewMenuParam(ctx)
	if param.Alert != "" {
		ctx.Header("Context-Type", "text/html; charset=utf-8")
		ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuURL)
		h.getMenuInfoPanel(ctx, param.Alert)
		return
	}
	// GetUserByMiddleware 取得middleware驗證後的user
	user := auth.GetUserByMiddleware()

	// 新建菜單
	menuModel, err := models.DefaultMenuModel().SetConn(h.Conn).
		New(param.Title, param.Icon, param.URL, param.Header, param.ParentID, (menu.GetMenuInformation(user, h.Conn)).MaxOrder+1)
	if err != nil {
		h.showNewMenu(ctx, "新建菜單發生錯誤")
		return
	}

	for _, roleID := range param.Roles {
		_, err = menuModel.AddRole(roleID)
		if err != nil {
			h.showNewMenu(ctx, "新建角色發生錯誤")
			return
		}
	}

	// 增加MaxOrder
	menu.GetMenuInformation(user, h.Conn).MaxOrder++

	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuURL)
	h.getMenuInfoPanel(ctx, "")
}

// ShowNewMenu new menu GET功能
func (h *Handler) ShowNewMenu(ctx *gin.Context) {
	h.showNewMenu(ctx, "")
}

// showNewMenu new menu面板
func (h *Handler) showNewMenu(ctx *gin.Context, alert string) {
	// GetUserByMiddleware 取得middleware驗證後的user
	user := auth.GetUserByMiddleware()
	// GetMenuInformation 透過user取得menu資料表資訊
	menuInfo := menu.GetMenuInformation(user, h.Conn)
	// 取得表單資訊
	formInfo := table.GetMenuFormPanel(h.Conn).GetNewForm(h.Services)
	formID, err := uuid.NewV4()
	if err != nil {
		panic("產生uuid發生錯誤")
	}
	urlPrefix := "/" + h.Config.URLPrefix

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
		MiniLogo     template.HTML
		Logo         template.HTML
		Previous     string
		Token        string
		URL          string
		URLPrefix    string
		IndexURL     string
	}{
		FormID:       formID.String(),
		User:         user,
		Content:      formInfo.FieldList,
		AlertContent: alert,
		Menu:         menuInfo,
		MiniLogo:     h.Config.MiniLogo,
		Logo:         h.Config.Logo,
		Previous:     urlPrefix + h.Config.MenuURL,
		Token:        auth.ConvertInterfaceToTokenService(h.Services.Get("token_csrf_helper")).AddToken(),
		URL:          urlPrefix + h.Config.MenuNewURL,
		URLPrefix:    urlPrefix,
		IndexURL:     urlPrefix + h.Config.IndexURL,
	}); err == nil {
		ctx.Status(http.StatusOK)
		ctx.Header("Content-Type", "text/html; charset=utf-8")
	} else {
		ctx.Status(http.StatusOK)
		panic("使用新建菜單模板發生錯誤")
	}
}