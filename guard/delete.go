package guard

import (
	"hilive/context"
	"hilive/modules/auth"
	"hilive/modules/parameter"
	"hilive/modules/table"
	"net/http"
)

// DeleteParam 設置刪除POST功能參數
type DeleteParam struct {
	Panel  table.Table
	ID     string
	Prefix string
	Param  parameter.Parameters
}

// Delete 取得url的id值後將值設置至Context.UserValue[delete_param]
func (g *Guard) Delete(ctx *context.Context) {

	// GetTable 取得table(面板資訊、表單資訊)
	panel, prefix := g.GetTable(ctx)

	id := ctx.Request.FormValue("id")
	if id == "" {
		// 取得目前登入用戶(Context.UserValue["user"])並轉換成UserModel
		user := auth.Auth(ctx)

		buf := g.ExecuteAlert(ctx, user, "刪除資料需要設置id參數")
		ctx.HTML(http.StatusOK, buf.String())
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
func GetDeleteParam(ctx *context.Context) *DeleteParam {
	return parameters["delete_parameter"].(*DeleteParam)
}
