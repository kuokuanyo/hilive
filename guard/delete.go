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

// DeleteParam 設置刪除POST功能參數
type DeleteParam struct {
	Panel  table.Table
	ID     string
	Prefix string
	Param  parameter.Parameters
}

// Delete 取得url的id值後將值設置至Context.UserValue[delete_param]
func (g *Guard) Delete(ctx *gin.Context) {
	// GetTable 取得table(面板資訊、表單資訊)
	panel, prefix := g.GetTable(ctx)

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
			AlertContent: "刪除資料需要設置id參數",
			Config:       g.Config,
			URLPrefix:    config.Prefix(),
			IndexURL:     config.Prefix() + g.Config.IndexURL,
		})
		ctx.Abort()
		return
	}

	parameters["delete_parameter"] = &DeleteParam{
		Panel:  panel,
		ID:     id,
		Prefix: prefix,
		Param:  parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize),
	}
	ctx.Next()
}

// GetDeleteParam 取得parameters["delete_form_parameter"]
func GetDeleteParam(ctx *gin.Context) *DeleteParam {
	return parameters["delete_parameter"].(*DeleteParam)
}
