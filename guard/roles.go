package guard

import (
	"hilive/modules/parameter"
	"hilive/modules/table"

	"github.com/gin-gonic/gin"
)

// ShowRolesNewForm 將新增角色的資訊設置至show_roles_new_form_parameter
func (g *Guard) ShowRolesNewForm(ctx *gin.Context) {
	panel := table.GetRolesFormPanel(g.Conn)

	parameters["show_roles_new_form_parameter"] = &ShowNewFormParameter{
		Panel:  panel,
		Param:  parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize),
		Prefix: "roles",
	}
	ctx.Next()
}

// GetShowRolesNewForm 取得parameters["show_roles_new_form_parameter"]
func GetShowRolesNewForm(ctx *gin.Context) *ShowNewFormParameter {
	return parameters["show_roles_new_form_parameter"].(*ShowNewFormParameter)
}
