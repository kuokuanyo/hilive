package guard

import (
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/menu"
	"hilive/views/alert"
	"html/template"
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
		Roles:    ctx.Request.Form["roles[]"],
		Alert:    alert,
	}
	ctx.Next()
}

// MenuDelete 建立刪除菜單參數
func (g *Guard) MenuDelete(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		user := auth.GetUserByMiddleware()
		// GetMenuInformation 透過user取得menu資料表資訊
		menuInfo := menu.GetMenuInformation(user, g.Conn)

		tmpl, err := template.New("").Funcs(template.FuncMap{
			"isLinkURL": func(s string) bool {
				return (len(s) > 7 && s[:7] == "http://") || (len(s) > 8 && s[:8] == "https://")
			},
		}).Parse(alert.AlertTmpl)
		if err != nil {
			panic("使用alert模板發生錯誤")
		}
		tmpl.Execute(ctx.Writer, struct {
			User         models.UserModel
			Menu         *menu.Menu
			AlertContent string
			Config       *config.Config
			URLPrefix    string
			IndexURL     string
		}{
			User:         user,
			Menu:         menuInfo,
			AlertContent: "刪除菜單需要設置id參數",
			Config:       g.Config,
			URLPrefix:    config.Prefix(),
			IndexURL:     config.Prefix() + g.Config.IndexURL,
		})
		ctx.Abort()
		return
	}

	parameters["delete_menu_param"] = &MenuDeleteParam{
		ID: id,
	}
	ctx.Next()
}

// GetNewMenuParam 取得Parameters["new_menu_param"]
func GetNewMenuParam(ctx *gin.Context) *MenuNewParam {
	return parameters["new_menu_param"].(*MenuNewParam)
}

// GetEditMenuParam 取得parameters["edit_menu_param"]
func GetEditMenuParam(ctx *gin.Context) *MenuEditParam {
	return parameters["edit_menu_param"].(*MenuEditParam)
}

// GetDeleteMenuParam 取得parameters["delete_menu_param"]
func GetDeleteMenuParam(ctx *gin.Context) *MenuDeleteParam {
	return parameters["delete_menu_param"].(*MenuDeleteParam)
}
