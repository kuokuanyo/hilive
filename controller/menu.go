package controller

import (
	"fmt"
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/menu"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"hilive/template/types"
	"hilive/views/alert"
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

// ShowNewMenu new menu GET功能
func (h *Handler) ShowNewMenu(ctx *gin.Context) {
	h.showNewMenu(ctx, "")
}

// ShowEditMenu edit menu GET功能
func (h *Handler) ShowEditMenu(ctx *gin.Context) {
	if ctx.Query("id") == "" {
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.Header("X-Pjax-Url", "/"+h.Config.URLPrefix+h.Config.MenuURL)
		h.getMenuInfoPanel(ctx, "請填寫id參數")
		return
	}

	urlPrefix := "/" + h.Config.URLPrefix
	// 取得middleware驗證後的user
	user := auth.GetUserByMiddleware()
	// GetMenuInformation 透過user取得menu資料表資訊
	menuInfo := menu.GetMenuInformation(user, h.Conn)

	formInfo, err := table.GetMenuFormPanel(h.Conn).
		GetDataWithID(parameter.DefaultParameters().SetFieldPKByJoinParam(ctx.Query("id")), h.Services)

	if err != nil {
		tmpl, err := template.New("").Funcs(DefaultFuncMap).Parse(alert.AlertTmpl)
		if err != nil {
			panic("使用alert模板發生錯誤")
		}
		tmpl.Execute(ctx.Writer, struct {
			User         models.UserModel
			Menu         *menu.Menu
			AlertContent string
			MiniLogo     template.HTML
			Logo         template.HTML
			IndexURL     string
			URLPrefix    string
		}{
			User:         user,
			Menu:         menuInfo,
			AlertContent: err.Error(),
			MiniLogo:     h.Config.MiniLogo,
			Logo:         h.Config.Logo,
			IndexURL:     urlPrefix + h.Config.IndexURL,
			URLPrefix:    urlPrefix,
		})
	}

	formID, err := uuid.NewV4()
	if err != nil {
		panic("產生uuid發生錯誤")
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
		MiniLogo     template.HTML
		Logo         template.HTML
		URL          string
		IndexURL     string
		URLPrefix    string
		Previous     string
		Token        string
	}{
		FormID:       formID.String(),
		User:         user,
		Menu:         menuInfo,
		AlertContent: "",
		Content:      formInfo.FieldList,
		MiniLogo:     h.Config.MiniLogo,
		Logo:         h.Config.Logo,
		URL:          urlPrefix + h.Config.MenuEditURL,
		IndexURL:     urlPrefix + h.Config.IndexURL,
		URLPrefix:    urlPrefix,
		Previous:     urlPrefix + h.Config.MenuURL,
		Token:        auth.ConvertInterfaceToTokenService(h.Services.Get("token_csrf_helper")).AddToken(),
	})
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
		fmt.Println(err)
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
		AlertContent string
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
		AlertContent: alert,
		Content:      formInfo.FieldList,
		NewMenuTitle: "新建菜單",
		Title:        "菜單",
		Description:  "菜單處理",
		Previous:     urlPrefix + h.Config.MenuURL,
		Token:        auth.ConvertInterfaceToTokenService(h.Services.Get("token_csrf_helper")).AddToken(),
		MiniLogo:     h.Config.MiniLogo,
		Logo:         h.Config.Logo,
		URL:          urlPrefix + h.Config.MenuURL,
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
