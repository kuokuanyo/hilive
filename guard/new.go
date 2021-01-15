package guard

import (
	"hilive/context"
	"hilive/modules/auth"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"html/template"
	"mime/multipart"
	"net/http"
	"strings"
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
func (g *Guard) NewForm(ctx *context.Context) {
	// GetTable 取得table(面板資訊、表單資訊)
	panel, prefix := g.GetTable(ctx)
	previous := ctx.Request.FormValue("__previous_")

	// 檢查token
	token := ctx.Request.FormValue("__token_")
	if !auth.GetTokenServiceByService(g.services.Get("token_csrf_helper")).CheckToken(token) {
		// 取得目前登入用戶(Context.UserValue["user"])並轉換成UserModel
		user := auth.Auth(ctx)

		buf := g.ExecuteAlert(ctx, user, "錯誤的token")
		ctx.HTML(http.StatusOK, buf.String())
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
func GetNewForm(ctx *context.Context) *NewFormParameter {
	return parameters["new_form_parameter"].(*NewFormParameter)
}
