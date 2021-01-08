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
	"mime/multipart"
	"strings"

	"github.com/gin-gonic/gin"
)

// NewFormParameter 設置新增資料POST功能資訊
type NewFormParameter struct {
	Panel     table.Table
	ID        string
	Prefix    string
	Param     parameter.Parameters
	Path      string
	MultiForm *multipart.Form
	Alert     template.HTML
}

// NewForm 設置新增用戶POST 功能資訊至new_form_parameter
func (g *Guard) NewForm(ctx *gin.Context) {
	// GetTable 取得table(面板資訊、表單資訊)
	panel, prefix := g.GetTable(ctx)
	previous := ctx.Request.FormValue("__previous_")

	// 檢查token
	token := ctx.Request.FormValue("__token_")
	if !auth.GetTokenServiceByService(g.Services.Get("token_csrf_helper")).CheckToken(token) {
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
			AlertContent: "錯誤的token",
			Config:       g.Config,
			URLPrefix:    config.Prefix(),
			IndexURL:     config.Prefix() + g.Config.IndexURL,
		})
		ctx.Abort()
		return
	}

	// GetParamFromURL 解析URL後設置頁面資訊
	param := parameter.GetParamFromURL(previous, panel.GetInfo().DefaultPageSize)

	parameters["new_form_parameter"] = &NewFormParameter{
		Panel:     panel,
		ID:        "",
		Prefix:    prefix,
		Param:     param,
		MultiForm: ctx.Request.MultipartForm,
		Path:      strings.Split(previous, "?")[0],
	}
	ctx.Next()
}

// GetNewForm 取得parameters["new_form_parameter"]
func GetNewForm(ctx *gin.Context) *NewFormParameter {
	return parameters["new_form_parameter"].(*NewFormParameter)
}
