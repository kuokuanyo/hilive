package guard

import (
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/menu"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"hilive/views/alert"
	"html/template"

	"github.com/gin-gonic/gin"
)

// ShowEditFormParameter 設置編輯表單POST功能資訊
type ShowEditFormParameter struct {
	Panel  table.Table
	ID     string
	Prefix string
	Param  parameter.Parameters
}

// ShowNewFormParameter 設置新增表單GET功能資訊
type ShowNewFormParameter struct {
	Panel  table.Table          // 表單資訊
	Param  parameter.Parameters // 頁面資訊
	Prefix string
}

// ShowNewForm 將新增用戶GET功能資訊設置至show_new_form_parameter
func (g *Guard) ShowNewForm(ctx *gin.Context) {
	// GetTable 取得table(面板資訊、表單資訊)
	panel, prefix := g.GetTable(ctx)

	parameters["show_new_form_parameter"] = &ShowNewFormParameter{
		Panel:  panel,
		Param:  parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize),
		Prefix: prefix,
	}
	ctx.Next()
}

// ShowEditForm 將編輯資訊設置至show_edit_form_parameter
func (g *Guard) ShowEditForm(ctx *gin.Context) {
	// GetTable 取得table(面板資訊、表單資訊)
	panel, prefix := g.GetTable(ctx)

	id := ctx.Query("__edit_pk")
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
			AlertContent: "編輯功能需要設置__edit_pk參數",
			Config:       g.Config,
			URLPrefix:    config.Prefix(),
			IndexURL:     config.Prefix() + g.Config.IndexURL,
		})
		ctx.Abort()
		return
	}
	parameters["show_edit_form_parameter"] = &ShowEditFormParameter{
		Panel:  panel,
		ID:     id,
		Prefix: prefix,
		Param:  parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize).SetPKs(id),
	}
	ctx.Next()
}

// GetShowNewForm 取得parameters["show_new_form_parameter"]
func GetShowNewForm(ctx *gin.Context) *ShowNewFormParameter {
	return parameters["show_new_form_parameter"].(*ShowNewFormParameter)
}

// GetShowEditForm 取得parameters["show_edit_form_parameter"]
func GetShowEditForm(ctx *gin.Context) *ShowEditFormParameter {
	return parameters["show_edit_form_parameter"].(*ShowEditFormParameter)
}
