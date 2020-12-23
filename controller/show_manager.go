package controller

import (
	"fmt"
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/menu"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"hilive/views/alert"
	"html/template"

	"github.com/gin-gonic/gin"
)

// ShowManegerInfo 前端資訊頁面
func (h *Handler) ShowManegerInfo(ctx *gin.Context) {
	// GetManagerInfoPanel 取得使用者資訊面板
	panel := table.GetManagerInfoPanel(h.Conn)
	// GetParam 設置頁面資訊
	params := parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize)
	h.showManagerTable(ctx, params, panel)
}

// showManagerTable 取得頁面資料後並執行前端模板語法
func (h *Handler) showManagerTable(ctx *gin.Context, params parameter.Parameters, panel table.Table) {
	// 取得middleware驗證後的user
	user := auth.GetUserByMiddleware()
	// GetMenuInformation 透過user取得menu資料表資訊
	menuInfo := menu.GetMenuInformation(user, h.Conn)
	urlPrefix := "/" + h.Config.URLPrefix

	panel, _, urls, err := h.ManegerTableData(ctx, params, panel)
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
			AlertContent: "取得頁面需要使用的資訊發生錯誤",
			MiniLogo:     h.Config.MiniLogo,
			Logo:         h.Config.Logo,
			IndexURL:     urlPrefix + h.Config.IndexURL,
			URLPrefix:    urlPrefix,
		})
		return
	}

	editURL, newURL, deleteURL, infoURL := urls[0], urls[1], urls[2], urls[3]
	fmt.Println(editURL)
	fmt.Println(newURL)
	fmt.Println(deleteURL)
	fmt.Println(infoURL)
}

// ManegerTableData 取得前端頁面需要的資訊(每一筆資料資訊、欄位資訊、可過濾欄位資訊...等)，並檢查是否有權限訪問URL
func (h *Handler) ManegerTableData(ctx *gin.Context, params parameter.Parameters,
	panel table.Table) (table.Table, table.PanelInfo, []string, error) {
	if panel == nil {
		table.GetManagerInfoPanel(h.Conn)
	}

	panelInfo, err := panel.GetData(params, h.Services)
	if err != nil {
		return panel, panelInfo, nil, err
	}

	// url後的參數(ex: ?__page=1&__pageSize=10&__sort=id&__sort_type=desc)
	paramStr := params.GetRouteParamStr()
	// 按鈕需要使用的URL
	editURL := "/" + h.Config.URLPrefix + h.Config.ManagerURL + "/edit" + paramStr
	newURL := "/" + h.Config.URLPrefix + h.Config.ManagerURL + "/new" + paramStr
	deleteURL := "/" + h.Config.URLPrefix + "/delete/manager"
	infoURL := "/" + h.Config.URLPrefix + h.Config.ManagerURL

	// 檢查用戶權限，如果沒有權限則回傳空的URL
	user := auth.GetUserByMiddleware()
	editURL = user.GetCheckPermissionByURLMethod(editURL, "GET")
	newURL = user.GetCheckPermissionByURLMethod(newURL, "GET")
	deleteURL = user.GetCheckPermissionByURLMethod(deleteURL, "POST")

	return panel, panelInfo, []string{editURL, newURL, deleteURL, infoURL}, nil
}
