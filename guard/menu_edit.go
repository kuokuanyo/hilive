package guard

import (
	"hilive/modules/auth"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

// MenuEdit 建立編輯菜單參數
func (g *Guard) MenuEdit(ctx *gin.Context) {
	parentID := ctx.Request.FormValue("parent_id")
	if parentID == "" {
		parentID = "0"
	}

	var (
		parentIDInt, _ = strconv.Atoi(parentID)
		token          = ctx.Request.FormValue("__token_")
		alert          string
	)

	if !auth.GetTokenServiceByService(g.Services.Get("token_csrf_helper")).CheckToken(token) {
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
		Roles:    ctx.Request.Form["roles"],
		Alert:    alert,
	}
	ctx.Next()
}

// GetEditMenuParam 取得parameters["edit_menu_param"]
func GetEditMenuParam(ctx *gin.Context) *MenuEditParam {
	return parameters["edit_menu_param"].(*MenuEditParam)
}
