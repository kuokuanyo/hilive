package guard

import (
	"hilive/context"
	"hilive/modules/auth"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"net/http"
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
func (g *Guard) ShowNewForm(ctx *context.Context) {
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
func (g *Guard) ShowEditForm(ctx *context.Context) {
	// GetTable 取得table(面板資訊、表單資訊)
	panel, prefix := g.GetTable(ctx)

	id := ctx.Query("__edit_pk")
	if id == "" {
		// 取得目前登入用戶(Context.UserValue["user"])並轉換成UserModel
		user := auth.Auth(ctx)

		buf := g.ExecuteAlert(ctx, user, "編輯功能需要設置__edit_pk參數")
		ctx.HTML(http.StatusOK, buf.String())
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
func GetShowNewForm(ctx *context.Context) *ShowNewFormParameter {
	return parameters["show_new_form_parameter"].(*ShowNewFormParameter)
}

// GetShowEditForm 取得parameters["show_edit_form_parameter"]
func GetShowEditForm(ctx *context.Context) *ShowEditFormParameter {
	return parameters["show_edit_form_parameter"].(*ShowEditFormParameter)
}
