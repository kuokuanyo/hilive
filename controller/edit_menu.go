package controller

import (
	"hilive/guard"
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/menu"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"hilive/template/types"
	"hilive/views/alert"
	"hilive/views/menuviews"
	"html/template"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// EditMenu 編輯菜單POST功能
func (h *Handler) EditMenu(ctx *gin.Context) {
	param := guard.GetEditMenuParam(ctx)
	if param.Alert != "" {
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuURL)
		h.getMenuInfoPanel(ctx, param.Alert)
		return
	}

	// 建立MenuModel
	menuModel := models.SetMenuModelByID(param.ID).SetConn(h.Conn)

	// 先刪除所有角色
	err := menuModel.DeleteRoles()
	if err != nil {
		if err.Error() != "沒有影響任何資料" {
			ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuURL)
			formInfo, _ := table.GetMenuFormPanel(h.Conn).
				GetDataWithID(parameter.DefaultParameters().SetFieldPKByJoinParam(param.ID), h.Services)
			h.showEditMenu(ctx, formInfo, "刪除角色發生錯誤")
			return
		}
	}

	// 新建角色
	for _, roleID := range param.Roles {
		_, err = menuModel.AddRole(roleID)
		if err != nil {
			ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuURL)
			formInfo, _ := table.GetMenuFormPanel(h.Conn).
				GetDataWithID(parameter.DefaultParameters().SetFieldPKByJoinParam(param.ID), h.Services)
			h.showEditMenu(ctx, formInfo, "新建角色發生錯誤")
			return
		}
	}

	// 更新資料
	_, err = menuModel.Update(param.Title, param.Icon, param.URL, param.Header, param.ParentID)
	if err != nil {
		if err.Error() != "沒有影響任何資料" {
			ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuURL)
			formInfo, _ := table.GetMenuFormPanel(h.Conn).
				GetDataWithID(parameter.DefaultParameters().SetFieldPKByJoinParam(param.ID), h.Services)
			h.showEditMenu(ctx, formInfo, "更新菜單資料發生錯誤")
			return
		}
	}

	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Header("X-PJAX-Url", "/"+h.Config.URLPrefix+h.Config.MenuURL)
	h.getMenuInfoPanel(ctx, "")
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

	h.showEditMenu(ctx, formInfo, "")
}

func (h *Handler) showEditMenu(ctx *gin.Context, formInfo table.FormInfo, alert string) {
	formID, err := uuid.NewV4()
	if err != nil {
		panic("產生uuid發生錯誤")
	}

	// 取得middleware驗證後的user
	user := auth.GetUserByMiddleware()
	// GetMenuInformation 透過user取得menu資料表資訊
	menuInfo := menu.GetMenuInformation(user, h.Conn)
	urlPrefix := "/" + h.Config.URLPrefix

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
		AlertContent: alert,
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
