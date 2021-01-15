package controller

import (
	"hilive/context"
	"hilive/guard"
	"net/http"
)

// EditForm 更新資料POST功能
func (h *Handler) EditForm(ctx *context.Context) {
	param := guard.GetEditForm(ctx)

	err := param.Panel.UpdateData(param.MultiForm.Value)
	if err != nil {
		h.showEditForm(ctx, err.Error(), param.Panel, param.Param, param.Prefix)
		ctx.AddHeader("X-PJAX-Url", param.Path+param.Param.DeletePK().DeleteEditPk().GetRouteParamStr())
		return
	}

	buf := h.showTable(ctx, param.Param.DeletePK(), param.Panel, param.Prefix, "")
	ctx.HTML(http.StatusOK, buf.String())
	ctx.AddHeader("X-PJAX-Url", param.Path+param.Param.DeletePK().DeleteEditPk().GetRouteParamStr())
}
