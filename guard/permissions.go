package guard

import (
	"hilive/modules/parameter"
	"hilive/modules/table"

	"github.com/gin-gonic/gin"
)

// ShowPermissionNewForm 將新增權限的資訊設置至show_permission_new_form_parameter
func (g *Guard) ShowPermissionNewForm(ctx *gin.Context) {
	panel := table.GetPermissionFormPanel(g.Conn)

	parameters["show_permission_new_form_parameter"] = &ShowNewFormParameter{
		Panel:  panel,
		Param:  parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize),
		Prefix: "permission",
	}
	ctx.Next()
}

// GetShowPermissionNewForm 取得parameters["show_permission_new_form_parameter"]
func GetShowPermissionNewForm(ctx *gin.Context) *ShowNewFormParameter {
	return parameters["show_permission_new_form_parameter"].(*ShowNewFormParameter)
}