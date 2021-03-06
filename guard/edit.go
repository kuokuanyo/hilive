package guard

import (
	"hilive/context"
	"hilive/modules/auth"
	"hilive/modules/parameter"
	"net/http"
	"strings"
)

// EditForm 編輯表單POST功能
func (g *Guard) EditForm(ctx *context.Context) {
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

	// 取得在multipart/form-data所設定的參數(struct)
	multiForm := ctx.Request.MultipartForm
	// 取得id
	id := multiForm.Value[panel.GetPrimaryKey().Name][0]

	// GetParamFromURL 解析URL後設置頁面資訊
	param := parameter.GetParamFromURL(previous, panel.GetInfo().DefaultPageSize)

	parameters["edit_form_parameter"] = &NewFormParameter{
		Panel:     panel,
		ID:        id,
		Prefix:    prefix,
		Param:     param.SetPKs(id),
		MultiForm: multiForm,
		Path:      strings.Split(previous, "?")[0],
	}
	ctx.Next()
}

// GetEditForm 取得parameters["edit_form_parameter"]
func GetEditForm(ctx *context.Context) *NewFormParameter {
	return parameters["edit_form_parameter"].(*NewFormParameter)
}
