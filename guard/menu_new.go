package guard

import (
	"hilive/modules/auth"
	"strconv"

	"github.com/gin-gonic/gin"
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

// MenuNew 建立Parameters[new_menu_param]
func (g *Guard) MenuNew(ctx *gin.Context) {
	var (
		token = ctx.Request.FormValue("__token_")
		alert string
	)

	if !auth.GetTokenServiceByService(g.Services.Get("token_csrf_helper")).CheckToken(token) {
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

// GetNewMenuParam 取得Parameters["new_menu_param"]
func GetNewMenuParam(ctx *gin.Context) *MenuNewParam {
	return parameters["new_menu_param"].(*MenuNewParam)
}
