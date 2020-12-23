package guard

import (
	"hilive/models"
	"hilive/modules/auth"
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
		urlPrefix := "/" + g.Config.URLPrefix
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
			MiniLogo     template.HTML
			Logo         template.HTML
			IndexURL     string
			URLPrefix    string
		}{
			User:         user,
			Menu:         menuInfo,
			AlertContent: "刪除菜單需要設置id參數",
			MiniLogo:     g.Config.MiniLogo,
			Logo:         g.Config.Logo,
			IndexURL:     urlPrefix + g.Config.IndexURL,
			URLPrefix:    urlPrefix,
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
