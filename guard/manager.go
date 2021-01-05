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

// ShowNewFormParameter 設置新增頁面的表單資訊及頁面資訊
type ShowNewFormParameter struct {
	Panel  table.Table          // 表單資訊
	Param  parameter.Parameters // 頁面資訊
	Prefix string
}

// ShowManagerForm 將編輯用戶資訊設置至show_manager_form_parameter
func (g *Guard) ShowManagerForm(ctx *gin.Context) {
	// GetManagerFormPanel 取得用戶表單資訊
	_ = table.GetManagerFormPanel(g.Conn)

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
			AlertContent: "編輯用戶需要設置__edit_pk參數",
			Config:       g.Config,
			URLPrefix:    config.Prefix(),
			IndexURL:     config.Prefix() + g.Config.IndexURL,
		})
		ctx.Abort()
		return
	}
}

// ShowManagerNewForm 將新增用戶的資訊設置至show_manager_new_form_parameter
func (g *Guard) ShowManagerNewForm(ctx *gin.Context) {
	panel := table.GetManagerFormPanel(g.Conn)

	parameters["show_manager_new_form_parameter"] = &ShowNewFormParameter{
		Panel:  panel,
		Param:  parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize),
		Prefix: "manager",
	}
	ctx.Next()
}

// GetShowManagerNewForm 取得parameters["show_manager_new_form_parameter"]
func GetShowManagerNewForm(ctx *gin.Context) *ShowNewFormParameter {
	return parameters["show_manager_new_form_parameter"].(*ShowNewFormParameter)
}
