package controller

import (
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/menu"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"hilive/modules/utils"
	"hilive/views/alert"
	"hilive/views/info"
	"html/template"

	"github.com/gin-gonic/gin"
)

// ShowManegerInfo 前端資訊頁面
func (h *Handler) ShowManegerInfo(ctx *gin.Context) {
	// GetManagerInfoPanel 取得使用者資訊面板
	panel := table.GetManagerInfoPanel(h.Conn)
	// GetParam 設置頁面資訊
	params := parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize)

	// 取得頁面資料後並執行前端模板語法
	h.showTable(ctx, params, panel, h.Config.ManagerURL, "manager")
}

// showManagerTable 取得頁面資料後並執行前端模板語法
func (h *Handler) showTable(ctx *gin.Context, params parameter.Parameters, panel table.Table, url string, prefix string) {
	// 取得middleware驗證後的user
	user := auth.GetUserByMiddleware()
	// GetMenuInformation 透過user取得menu資料表資訊
	menuInfo := menu.GetMenuInformation(user, h.Conn).SetActiveClass(url)

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
