package controller

import (
	"bytes"
	"hilive/context"
	"hilive/guard"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"net/http"
)

// ShowEditForm 編輯資料前端頁面
func (h *Handler) ShowEditForm(ctx *context.Context) {
	param := guard.GetShowEditForm(ctx)

	h.showEditForm(ctx, "", param.Panel, param.Param, param.Prefix)
}

// ShowInfo 前端資訊頁面
func (h *Handler) ShowInfo(ctx *context.Context) {
	prefix := ctx.Query("__prefix")

	// GetTable 取得table(面板資訊、表單資訊)
	panel := h.GetTable(ctx, prefix)
	// GetParam 設置頁面資訊
	params := parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize)

	// 取得頁面資料後並執行前端模板語法
	buf := h.showTable(ctx, params, panel, prefix, "")
	ctx.HTML(http.StatusOK, buf.String())
}

// ShowNewForm 新增用戶前端頁面
func (h *Handler) ShowNewForm(ctx *context.Context) {
	param := guard.GetShowNewForm(ctx)
	h.showNewForm(ctx, "", param.Panel, param.Param.GetRouteParamStr(), param.Prefix)
}

// showEditForm 編輯頁面模板語法
func (h *Handler) showEditForm(ctx *context.Context, alertMsg string, panel table.Table, param parameter.Parameters, prefix string) {
	// 取得目前登入用戶(Context.UserValue["user"])並轉換成UserModel
	user := auth.Auth(ctx)

	// GetDataWithID 透過id取得資料並將值、預設值設置至BaseTable.Form.FormFields
	formInfo, err := panel.GetDataWithID(param, h.services)
	if err != nil {
		buf := h.ExecuteAlert(ctx, user, err.Error())
		ctx.HTML(http.StatusOK, buf.String())
		return
	}

	route := URLRoute{
		URLPrefix:   config.Prefix(),
		InfoURL:     config.Prefix() + "/edit/" + prefix,
		IndexURL:    config.Prefix() + h.config.IndexURL,
		PreviousURL: config.Prefix() + "/info/" + prefix + param.DeletePK().DeleteField("__edit_pk").GetRouteParamStr(),
	}

	buf := h.ExecuteForm(ctx, user, formInfo, route, alertMsg, "form")
	ctx.HTML(http.StatusOK, buf.String())
}

// showNewForm 新增功能模板語法
func (h *Handler) showNewForm(ctx *context.Context, alert string, panel table.Table, paramStr string, prefix string) {
	// 取得目前登入用戶(Context.UserValue["user"])並轉換成UserModel
	user := auth.Auth(ctx)

	// GetNewForm 處理並設置表單欄位細節資訊(允許增加的表單欄位)
	formInfo := panel.GetNewForm(h.services)

	route := URLRoute{
		URLPrefix:   config.Prefix(),
		InfoURL:     config.Prefix() + "/new/" + prefix,
		IndexURL:    config.Prefix() + h.config.IndexURL,
		PreviousURL: config.Prefix() + "/info/" + prefix + paramStr,
	}

	buf := h.ExecuteForm(ctx, user, formInfo, route, alert, "form")
	ctx.HTML(http.StatusOK, buf.String())
}

// showTable 取得頁面資料後並執行前端模板語法
func (h *Handler) showTable(ctx *context.Context, params parameter.Parameters, panel table.Table, prefix, alertmsg string) *bytes.Buffer {
	// 取得目前登入用戶(Context.UserValue["user"])並轉換成UserModel
	user := auth.Auth(ctx)

	// GetTableData 取得前端頁面需要的資訊(每一筆資料資訊、欄位資訊、可過濾欄位資訊...等)，並檢查是否有權限訪問URL
	panel, panelInfo, urls, err := h.GetTableData(ctx, params, panel, prefix)
	if err != nil {
		return h.ExecuteAlert(ctx, user, "取得頁面需要使用的資訊發生錯誤")
	}

	// 設置模板需要用到的url路徑
	editURL, newURL, deleteURL, infoURL := urls[0], urls[1], urls[2], urls[3]
	route := URLRoute{
		IndexURL:  config.Prefix() + h.config.IndexURL,
		URLPrefix: config.Prefix(),
		InfoURL:   infoURL,
		SortURL:   params.GetFixedParamWithoutSort(),
		NewURL:    newURL,
		EditURL:   editURL,
		DeleteURL: deleteURL,
	}

	return h.ExecuteInfo(ctx, user, panelInfo, route, "")
}

// GetTableData 取得前端頁面需要的資訊(每一筆資料資訊、欄位資訊、可過濾欄位資訊...等)，並檢查是否有權限訪問URL
func (h *Handler) GetTableData(ctx *context.Context, params parameter.Parameters,
	panel table.Table, prefix string) (table.Table, table.PanelInfo, []string, error) {
	// 先設置table(interface)
	if panel == nil {
		panel = h.GetTable(ctx, prefix)
	}

	// 從資料庫取得頁面需要顯示的資料，回傳每一筆資料資訊、欄位資訊、可過濾欄位資訊、分頁資訊...等
	panelInfo, err := panel.GetData(params, h.services)
	if err != nil {
		return panel, panelInfo, nil, err
	}

	url := "/info/" + prefix
	// url後的參數(ex: ?__page=1&__pageSize=10&__sort=id&__sort_type=desc)
	paramStr := params.GetRouteParamStr()
	// 按鈕需要使用的URL
	editURL := config.Prefix() + url + "/edit" + paramStr
	newURL := config.Prefix() + url + "/new" + paramStr
	deleteURL := config.Prefix() + "/delete/" + prefix
	infoURL := config.Prefix() + url

	// 取得目前登入用戶(Context.UserValue["user"])並轉換成UserModel
	user := auth.Auth(ctx)
	editURL = user.GetCheckPermissionByURLMethod(editURL, "GET")
	newURL = user.GetCheckPermissionByURLMethod(newURL, "GET")
	deleteURL = user.GetCheckPermissionByURLMethod(deleteURL, "POST")

	return panel, panelInfo, []string{editURL, newURL, deleteURL, infoURL}, nil
}
