package guard

import (
	"hilive/context"
	"hilive/modules/auth"
	"net/http"
	"strconv"
)

// MenuNewParam 新建菜單參數
type MenuNewParam struct {
	ID       string
	Title    string
	Header   string
	ParentID int64
	Icon     string
	URL      string
	Roles    []string
	Alert    string
}

// MenuEditParam 編輯菜單的參數
type MenuEditParam struct {
	ID       string
	Title    string
	Header   string
	ParentID int64
	Icon     string
	URL      string
	Roles    []string
	Alert    string
}

// MenuDeleteParam 刪除菜單參數
type MenuDeleteParam struct {
	ID string
}

// MenuNew 建立Parameters[new_menu_param]
func (g *Guard) MenuNew(ctx *context.Context) {
	var (
		token = ctx.Request.FormValue("__token_")
		alert string
	)

	if !auth.GetTokenServiceByService(g.services.Get("token_csrf_helper")).CheckToken(token) {
		alert = "錯誤的token"
	}
	if alert == "" {
		if ctx.Request.FormValue("title") == "" || ctx.Request.FormValue("icon") == "" {
			alert = "title或icon參數不能為空"
		}
	}

	parentID := ctx.Request.FormValue("parent_id")
	if parentID == "" {
		parentID = "0"
	}
	parentIDInt, _ := strconv.Atoi(parentID)

	parameters["new_menu_param"] = &MenuNewParam{
		Title:    ctx.Request.FormValue("title"),
		Header:   ctx.Request.FormValue("header"),
		ParentID: int64(parentIDInt),
		Icon:     ctx.Request.FormValue("icon"),
		URL:      ctx.Request.FormValue("url"),
		Roles:    ctx.Request.Form["roles[]"],
		Alert:    alert,
	}
	ctx.Next()
}

// MenuEdit 建立編輯菜單參數
func (g *Guard) MenuEdit(ctx *context.Context) {
	parentID := ctx.Request.FormValue("parent_id")
	if parentID == "" {
		parentID = "0"
	}

	var (
		parentIDInt, _ = strconv.Atoi(parentID)
		token          = ctx.Request.FormValue("__token_")
		alert          string
	)

	if !auth.GetTokenServiceByService(g.services.Get("token_csrf_helper")).CheckToken(token) {
		alert = "錯誤的token"
	}
	if alert == "" {
		if ctx.Request.FormValue("id") == "" || ctx.Request.FormValue("title") == "" || ctx.Request.FormValue("icon") == "" {
			alert = "id、title、icon參數不能為空"
		}
	}

	parameters["edit_menu_param"] = &MenuEditParam{
		ID:       ctx.Request.FormValue("id"),
		Title:    ctx.Request.FormValue("title"),
		Header:   ctx.Request.FormValue("header"),
		ParentID: int64(parentIDInt),
		Icon:     ctx.Request.FormValue("icon"),
		URL:      ctx.Request.FormValue("url"),
		Roles:    ctx.Request.Form["roles[]"],
		Alert:    alert,
	}
	ctx.Next()
}

// MenuDelete 建立刪除菜單參數
func (g *Guard) MenuDelete(ctx *context.Context) {
	id := ctx.Query("id")
	if id == "" {
		// 取得目前登入用戶(Context.UserValue["user"])並轉換成UserModel
		user := auth.Auth(ctx)

		buf := g.ExecuteAlert(ctx, user, "刪除菜單需要設置id參數")
		ctx.HTML(http.StatusOK, buf.String())
		ctx.Abort()
		return
	}

	parameters["delete_menu_param"] = &MenuDeleteParam{
		ID: id,
	}
	ctx.Next()
}

// GetNewMenuParam 取得Parameters["new_menu_param"]
func GetNewMenuParam(ctx *context.Context) *MenuNewParam {
	return parameters["new_menu_param"].(*MenuNewParam)
}

// GetEditMenuParam 取得parameters["edit_menu_param"]
func GetEditMenuParam(ctx *context.Context) *MenuEditParam {
	return parameters["edit_menu_param"].(*MenuEditParam)
}

// GetDeleteMenuParam 取得parameters["delete_menu_param"]
func GetDeleteMenuParam(ctx *context.Context) *MenuDeleteParam {
	return parameters["delete_menu_param"].(*MenuDeleteParam)
}
