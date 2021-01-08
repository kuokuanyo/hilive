package guard

import (
	"hilive/models"
	"hilive/modules/auth"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/menu"
	"hilive/modules/service"
	"hilive/modules/table"
	"hilive/views/alert"
	"html/template"

	"github.com/gin-gonic/gin"
)

// Parameters 紀錄參數
var parameters = make(map[string]interface{})

// Guard struct
type Guard struct {
	Services   service.List
	Conn       db.Connection
	Config     *config.Config
	TablelList map[string]func(conn db.Connection) table.Table
}

// GetTable 取得table(面板資訊、表單資訊)
func (g *Guard) GetTable(ctx *gin.Context) (table.Table, string) {
	prefix := ctx.Param("__prefix")
	return g.TablelList[prefix](g.Conn), prefix
}

// CheckPrefix 檢查是否有__prefix頁面
func (g *Guard) CheckPrefix(ctx *gin.Context) {
	prefix := ctx.Param("__prefix")

	if _, ok := g.TablelList[prefix]; !ok {
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
			AlertContent: "抱歉，該頁面不存在!",
			Config:       g.Config,
			URLPrefix:    config.Prefix(),
			IndexURL:     config.Prefix() + g.Config.IndexURL,
		})
		ctx.Abort()
		return
	}
	ctx.Next()
}
