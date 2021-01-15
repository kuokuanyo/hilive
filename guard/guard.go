package guard

import (
	"bytes"
	"hilive/context"
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/menu"
	"hilive/modules/service"
	"hilive/modules/table"
	"hilive/views"
	"html/template"
	"net/http"
)

// Parameters 紀錄參數
var parameters = make(map[string]interface{})

// DefaultFuncMap 模板function
var DefaultFuncMap = template.FuncMap{
	"isLinkURL": func(s string) bool {
		return (len(s) > 7 && s[:7] == "http://") || (len(s) > 8 && s[:8] == "https://")
	},
}

// Guard struct
type Guard struct {
	services service.List
	conn     db.Connection
	config   *config.Config
	list     table.List
}

// ExecuteParam 使用模板需要的參數
type ExecuteParam struct {
	TmplName     string
	User         models.UserModel
	AlertContent string
	Menu         *menu.Menu
	Config       *config.Config
	URLRoute     URLRoute
}

// URLRoute 模板需要使用的URL路徑
type URLRoute struct {
	URLPrefix string
	IndexURL  string
}

// NewGuard 將Guard(struct)
func NewGuard(s service.List, c db.Connection, t table.List, cfg *config.Config) *Guard {
	return &Guard{
		services: s,
		conn:     c,
		list:     t,
		config:   cfg,
	}
}

// GetTable 取得table(面板資訊、表單資訊)
func (g *Guard) GetTable(ctx *context.Context) (table.Table, string) {
	prefix := ctx.Query("__prefix")
	return g.list[prefix](ctx), prefix
}

// CheckPrefix 檢查是否有__prefix頁面
func (g *Guard) CheckPrefix(ctx *context.Context) {
	prefix := ctx.Query("__prefix")

	if _, ok := g.list[prefix]; !ok {
		// 取得目前登入用戶(Context.UserValue["user"])並轉換成UserModel
		user := auth.Auth(ctx)

		buf := g.ExecuteAlert(ctx, user, "抱歉，該頁面不存在!")
		ctx.HTML(http.StatusOK, buf.String())
		ctx.Abort()
		return
	}
	ctx.Next()
}

// ExecuteAlert 執行錯誤警告模板
func (g *Guard) ExecuteAlert(ctx *context.Context, user models.UserModel, alertmsg string) *bytes.Buffer {
	var (
		tmpl *template.Template
		err  error
		name string
		buf  = new(bytes.Buffer)
	)

	if ctx.IsPjax() {
		name = "alert_content"
		tmpl, err = template.New(name).Funcs(DefaultFuncMap).Parse(views.TemplateList[name])
	} else {
		name = "layout"
		tmpl, err = template.New(name).Funcs(DefaultFuncMap).Parse(views.TemplateList[name] +
			views.TemplateList["head"] + views.TemplateList["header"] + views.TemplateList["sidebar"] +
			views.TemplateList["form_content"] + views.TemplateList["menu_content"] +
			views.TemplateList["info_content"] + views.TemplateList["alert_content"] +
			views.TemplateList["admin_panel"] + views.TemplateList["form"])
	}
	if err != nil {
		panic("使用錯誤警告模板發生錯誤")
	}
	err = tmpl.ExecuteTemplate(buf, name, ExecuteParam{
		TmplName:     "alert",
		User:         user,
		Menu:         menu.GetMenuInformation(user, g.conn).SetActiveClass(g.config.URLRemovePrefix(ctx.Path())),
		AlertContent: alertmsg,
		Config:       g.config,
		URLRoute: URLRoute{
			IndexURL:  config.Prefix() + g.config.IndexURL,
			URLPrefix: config.Prefix(),
		},
	})
	if err != nil {
		panic("執行錯誤警告模板發生錯誤")
	}
	return buf
}
