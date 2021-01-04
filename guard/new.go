package guard

import (
	"hilive/modules/parameter"
	"hilive/modules/table"

	"github.com/gin-gonic/gin"
)

// ShowNewFormParameter 設置新增頁面的表單資訊及頁面資訊
type ShowNewFormParameter struct {
	Panel  table.Table          // 表單資訊
	Param  parameter.Parameters // 頁面資訊
	Prefix string
}

// ShowManagerNewForm 將新增用戶的資訊設置至show_new_form_parameter
func (g *Guard) ShowManagerNewForm(ctx *gin.Context) {
	panel := table.GetManagerFormPanel(g.Conn)

	parameters["show_manager_new_form_parameter"] = &ShowNewFormParameter{
		Panel: panel,
		Param: parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize),
		Prefix: "manager",
	}
	ctx.Next()
}

// GetShowManagerNewForm 取得parameters["show_manager_new_form_parameter"]
func GetShowManagerNewForm(ctx *gin.Context) *ShowNewFormParameter {
	return parameters["show_manager_new_form_parameter"].(*ShowNewFormParameter)
}
