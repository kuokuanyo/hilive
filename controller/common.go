package controller

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
	"hilive/modules/utils"
	"hilive/views"
	"html/template"
	"regexp"
	"strings"
)

// Handler struct
type Handler struct {
	config   *config.Config
	conn     db.Connection
	services service.List
	list     table.List
	routes   context.RouterMap
}

// URLRoute 模板需要使用的URL路徑
type URLRoute struct {
	URLPrefix   string
	IndexURL    string
	InfoURL     string
	NewURL      string
	EditURL     string
	DeleteURL   string
	SortURL     string
	PreviousURL string
}

// ExecuteParam 使用模板需要的參數
type ExecuteParam struct {
	TmplName     string
	FormID       string
	User         models.UserModel
	PanelInfo    table.PanelInfo
	FormInfo     table.FormInfo
	AlertContent string
	Menu         *menu.Menu
	Config       *config.Config
	URLRoute     URLRoute
	Token        string
}

// ExecuteInfo 執行資訊面板的模板
func (h *Handler) ExecuteInfo(ctx *context.Context, user models.UserModel,
	panelInfo table.PanelInfo, route URLRoute, alertmsg string) *bytes.Buffer {
	var (
		tmpl *template.Template
		err  error
		name string
		buf  = new(bytes.Buffer)
	)

	if ctx.IsPjax() {
		name = "info_content"
		tmpl, err = template.New(name).Funcs(DefaultFuncMap).Parse(views.TemplateList[name] + views.TemplateList["form"])
	} else {
		name = "layout"
		tmpl, err = template.New(name).Funcs(DefaultFuncMap).Parse(views.TemplateList[name] +
			views.TemplateList["head"] + views.TemplateList["header"] + views.TemplateList["sidebar"] +
			views.TemplateList["form_content"] + views.TemplateList["menu_content"] +
			views.TemplateList["info_content"] + views.TemplateList["alert_content"] +
			views.TemplateList["admin_panel"] + views.TemplateList["form"])
	}
	if err != nil {
		panic("使用資訊面板模板發生錯誤")
	}
	err = tmpl.ExecuteTemplate(buf, name, ExecuteParam{
		TmplName:     "info",
		FormID:       utils.UUID(8),
		User:         user,
		PanelInfo:    panelInfo,
		AlertContent: alertmsg,
		Menu:         menu.GetMenuInformation(user, h.conn).SetActiveClass(h.config.URLRemovePrefix(ctx.Path())),
		Config:       h.config,
		URLRoute:     route,
	})
	if err != nil {
		panic("執行資訊面板模板發生錯誤")
	}
	return buf
}

// ExecuteForm 執行表單資訊的模板
func (h *Handler) ExecuteForm(ctx *context.Context, user models.UserModel,
	panelInfo table.FormInfo, route URLRoute, alertmsg string, TmplName string) *bytes.Buffer {
	var (
		tmpl *template.Template
		err  error
		name string
		buf  = new(bytes.Buffer)
	)

	if ctx.IsPjax() {
		name = TmplName + "_content"
		tmpl, err = template.New(name).Funcs(DefaultFuncMap).Parse(views.TemplateList[name] + views.TemplateList["form"])
	} else {
		name = "layout"
		tmpl, err = template.New(name).Funcs(DefaultFuncMap).Parse(views.TemplateList[name] +
			views.TemplateList["head"] + views.TemplateList["header"] + views.TemplateList["sidebar"] +
			views.TemplateList["form_content"] + views.TemplateList["menu_content"] +
			views.TemplateList["info_content"] + views.TemplateList["alert_content"] +
			views.TemplateList["form"] + views.TemplateList["admin_panel"])
	}
	if err != nil {
		panic("使用表單模板發生錯誤")
	}
	err = tmpl.ExecuteTemplate(buf, name, ExecuteParam{
		TmplName:     TmplName,
		FormID:       utils.UUID(8),
		User:         user,
		FormInfo:     panelInfo,
		AlertContent: alertmsg,
		Menu:         menu.GetMenuInformation(user, h.conn).SetActiveClass(h.config.URLRemovePrefix(ctx.Path())),
		Config:       h.config,
		URLRoute:     route,
		Token:        auth.ConvertInterfaceToTokenService(h.services.Get("token_csrf_helper")).AddToken(),
	})
	if err != nil {
		panic("執行表單模板發生錯誤")
	}
	return buf
}

// ExecuteAlert 執行錯誤警告模板
func (h *Handler) ExecuteAlert(ctx *context.Context, user models.UserModel, alertmsg string) *bytes.Buffer {
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
		Menu:         menu.GetMenuInformation(user, h.conn).SetActiveClass(h.config.URLRemovePrefix(ctx.Path())),
		AlertContent: alertmsg,
		Config:       h.config,
		URLRoute: URLRoute{
			IndexURL:  config.Prefix() + h.config.IndexURL,
			URLPrefix: config.Prefix(),
		},
	})
	if err != nil {
		panic("執行錯誤警告模板發生錯誤")
	}
	return buf
}

// NewHandler 設置Handler(struct)
func NewHandler() *Handler {
	return &Handler{}
}

// NewHandler 設置Handler(struct)
func (h *Handler) NewHandler(cfg *config.Config, services service.List, conn db.Connection, list table.List) {
	h.config = cfg
	h.services = services
	h.conn = conn
	h.list = list
}

// GetTable 取得table(面板資訊、表單資訊)
func (h *Handler) GetTable(ctx *context.Context, prefix string) table.Table {
	return h.list[prefix](ctx)
}

// DefaultFuncMap 模板需要使用的函式
var DefaultFuncMap = template.FuncMap{
	"link": func(cdnUrl, prefixUrl, assetsUrl string) string {
		if cdnUrl == "" {
			return prefixUrl + assetsUrl
		}
		return cdnUrl + assetsUrl
	},
	"isLinkURL": func(s string) bool {
		return (len(s) > 7 && s[:7] == "http://") || (len(s) > 8 && s[:8] == "https://")
	},
}

// isInfoURL 檢查url
func isInfoURL(s string) bool {
	reg, _ := regexp.Compile("(.*?)info/(.*?)$")
	sub := reg.FindStringSubmatch(s)
	return len(sub) > 2 && !strings.Contains(sub[2], "/")
}

// isNewURL 檢查url
func isNewURL(s string, p string) bool {
	reg, _ := regexp.Compile("(.*?)info/" + p + "/new")

	return reg.MatchString(s)
}

// isEditURL 檢查url
func isEditURL(s string, p string) bool {
	reg, _ := regexp.Compile("(.*?)info/" + p + "/edit")
	return reg.MatchString(s)
}
