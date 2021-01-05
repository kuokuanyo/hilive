package guard

import (
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/menu"
	"hilive/views/alert"
	"html/template"

	"github.com/gin-gonic/gin"
)

// MenuDeleteParam 刪除菜單參數
type MenuDeleteParam struct {
	ID string
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

// GetDeleteMenuParam 取得parameters["delete_menu_param"]
func GetDeleteMenuParam(ctx *gin.Context) *MenuDeleteParam {
	return parameters["delete_menu_param"].(*MenuDeleteParam)
}
