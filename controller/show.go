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
	"hilive/views/alert"
	"hilive/views/form"
	"hilive/views/info"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ShowEditForm 編輯資料前端頁面
func (h *Handler) ShowEditForm(ctx *gin.Context) {
	param := guard.GetShowEditForm(ctx)

	h.showEditForm(ctx, h.Alert, param.Panel, param.Param, param.Prefix)
}

// ShowInfo 前端資訊頁面
func (h *Handler) ShowInfo(ctx *gin.Context) {
	prefix := ctx.Param("__prefix")
	// GetTable 取得table(面板資訊、表單資訊
	panel := h.GetTable(ctx, prefix)
	// GetParam 設置頁面資訊
	params := parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize)

	// 取得頁面資料後並執行前端模板語法
	h.showTable(ctx, params, panel, "/info/"+prefix, prefix)
}

// ShowNewForm 新增用戶前端頁面
func (h *Handler) ShowNewForm(ctx *gin.Context) {
	param := guard.GetShowNewForm(ctx)
	h.showNewForm(ctx, h.Alert, param.Panel, param.Param.GetRouteParamStr(), param.Prefix)
}

// showEditForm 編輯頁面模板語法
func (h *Handler) showEditForm(ctx *gin.Context, alertMsg string, panel table.Table, param parameter.Parameters, prefix string) {
	// 取得middleware驗證後的user
	user := auth.GetUserByMiddleware()
	// GetMenuInformation 透過user取得menu資料表資訊
	menuInfo := menu.GetMenuInformation(user, h.Conn).SetActiveClass("/info/" + prefix + "/edit")

	// GetDataWithID 透過id取得資料並將值、預設值設置至BaseTable.Form.FormFields
	formInfo, err := panel.GetDataWithID(param, h.Services)
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

	route := URLRoute{
		URLPrefix:   config.Prefix(),
		InfoURL:     config.Prefix() + "/edit/" + prefix,
		IndexURL:    config.Prefix() + h.Config.IndexURL,
		PreviousURL: config.Prefix() + "/info/" + prefix + param.DeletePK().DeleteField("__edit_pk").GetRouteParamStr(),
	}

	tmpl, err := template.New("").Funcs(DefaultFuncMap).Parse(form.FormTmpl)
	if err != nil {
		panic("使用編輯模板發生錯誤")
	}
	if err := tmpl.Execute(ctx.Writer, struct {
		FormID       string
		User         models.UserModel
		FormInfo     table.FormInfo
		AlertContent string
		Menu         *menu.Menu
		Config       *config.Config
		URLRoute     URLRoute
		Token        string
	}{
		FormID:       utils.UUID(8),
		User:         user,
		FormInfo:     formInfo,
		AlertContent: alertMsg,
		Menu:         menuInfo,
		Config:       h.Config,
		URLRoute:     route,
		Token:        auth.ConvertInterfaceToTokenService(h.Services.Get("token_csrf_helper")).AddToken(),
	}); err == nil {
		ctx.Status(http.StatusOK)
		ctx.Header("Content-Type", "text/html; charset=utf-8")
	} else {
		panic("使用編輯用戶、角色、權限模板發生錯誤")
	}
	if alertMsg != "" {
		h.Alert = ""
	}
}

// showNewForm 新增功能模板語法
func (h *Handler) showNewForm(ctx *gin.Context, alert string, panel table.Table, paramStr string, prefix string) {
	// 取得middleware驗證後的user
	user := auth.GetUserByMiddleware()
	// GetMenuInformation 透過user取得menu資料表資訊
	menuInfo := menu.GetMenuInformation(user, h.Conn).SetActiveClass("/info/" + prefix + "/new")

	// GetNewForm 處理並設置表單欄位細節資訊(允許增加的表單欄位)
	formInfo := panel.GetNewForm(h.Services)

	route := URLRoute{
		URLPrefix:   config.Prefix(),
		InfoURL:     config.Prefix() + "/new/" + prefix,
		IndexURL:    config.Prefix() + h.Config.IndexURL,
		PreviousURL: config.Prefix() + "/info/" + prefix + paramStr,
	}

	tmpl, err := template.New("").Funcs(DefaultFuncMap).Parse(form.FormTmpl)
	if err != nil {
		panic("使用新建用戶、角色、權限模板發生錯誤")
	}
	if err := tmpl.Execute(ctx.Writer, struct {
		FormID       string
		User         models.UserModel
		FormInfo     table.FormInfo
		AlertContent string
		Menu         *menu.Menu
		Config       *config.Config
		URLRoute     URLRoute
		Token        string
	}{
		FormID:       utils.UUID(8),
		User:         user,
		FormInfo:     formInfo,
		AlertContent: alert,
		Menu:         menuInfo,
		Config:       h.Config,
		URLRoute:     route,
		Token:        auth.ConvertInterfaceToTokenService(h.Services.Get("token_csrf_helper")).AddToken(),
	}); err == nil {
		ctx.Status(http.StatusOK)
		ctx.Header("Content-Type", "text/html; charset=utf-8")
	} else {
		panic("使用新建用戶、角色、權限模板發生錯誤")
	}
	if alert != "" {
		h.Alert = ""
	}
}

// showTable 取得頁面資料後並執行前端模板語法
func (h *Handler) showTable(ctx *gin.Context, params parameter.Parameters, panel table.Table, url string, prefix string) {
	// 取得middleware驗證後的user
	user := auth.GetUserByMiddleware()
	// GetMenuInformation 透過user取得menu資料表資訊
	menuInfo := menu.GetMenuInformation(user, h.Conn).SetActiveClass(url)

	// GetTableData 取得前端頁面需要的資訊(每一筆資料資訊、欄位資訊、可過濾欄位資訊...等)，並檢查是否有權限訪問URL
	panel, panelInfo, urls, err := h.GetTableData(ctx, params, panel, url, prefix)
	if err != nil {
		route := URLRoute{
			IndexURL:  config.Prefix() + h.Config.IndexURL,
			URLPrefix: config.Prefix(),
		}
		tmpl, err := template.New("").Funcs(DefaultFuncMap).Parse(alert.AlertTmpl)
		if err != nil {
			panic("使用alert模板發生錯誤")
		}
		tmpl.Execute(ctx.Writer, struct {
			User         models.UserModel
			Menu         *menu.Menu
			AlertContent string
			Config       *config.Config
			URLRoute     URLRoute
		}{
			User:         user,
			Menu:         menuInfo,
			AlertContent: "取得頁面需要使用的資訊發生錯誤",
			Config:       h.Config,
			URLRoute:     route,
		})
		return
	}

	// 設置模板需要用到的url路徑
	editURL, newURL, deleteURL, infoURL := urls[0], urls[1], urls[2], urls[3]
	route := URLRoute{
		IndexURL:  config.Prefix() + h.Config.IndexURL,
		URLPrefix: config.Prefix(),
		InfoURL:   infoURL,
		SortURL:   params.GetFixedParamWithoutSort(),
		NewURL:    newURL,
		EditURL:   editURL,
		DeleteURL: deleteURL,
	}

	tmpl, err := template.New("").Funcs(DefaultFuncMap).Parse(info.InfoTmpl)
	if err != nil {
		panic("使用使用者介面模板發生錯誤")
	}
	tmpl.Execute(ctx.Writer, struct {
		FormID    string
		User      models.UserModel
		PanelInfo table.PanelInfo
		Menu      *menu.Menu
		Config    *config.Config
		URLRoute  URLRoute
	}{
		FormID:    utils.UUID(8),
		User:      user,
		PanelInfo: panelInfo,
		Menu:      menuInfo,
		Config:    h.Config,
		URLRoute:  route,
	})
}

// GetTableData 取得前端頁面需要的資訊(每一筆資料資訊、欄位資訊、可過濾欄位資訊...等)，並檢查是否有權限訪問URL
func (h *Handler) GetTableData(ctx *gin.Context, params parameter.Parameters,
	panel table.Table, url string, prefix string) (table.Table, table.PanelInfo, []string, error) {

	// 從資料庫取得頁面需要顯示的資料，回傳每一筆資料資訊、欄位資訊、可過濾欄位資訊、分頁資訊...等
	panelInfo, err := panel.GetData(params, h.Services)
	if err != nil {
		return panel, panelInfo, nil, err
	}

	// url後的參數(ex: ?__page=1&__pageSize=10&__sort=id&__sort_type=desc)
	paramStr := params.GetRouteParamStr()
	// 按鈕需要使用的URL
	editURL := config.Prefix() + url + "/edit" + paramStr
	newURL := config.Prefix() + url + "/new" + paramStr
	deleteURL := config.Prefix() + "/delete/" + prefix
	infoURL := config.Prefix() + url

	// 檢查用戶權限，如果沒有權限則回傳空的URL
	user := auth.GetUserByMiddleware()
	editURL = user.GetCheckPermissionByURLMethod(editURL, "GET")
	newURL = user.GetCheckPermissionByURLMethod(newURL, "GET")
	deleteURL = user.GetCheckPermissionByURLMethod(deleteURL, "POST")

	return panel, panelInfo, []string{editURL, newURL, deleteURL, infoURL}, nil
}
