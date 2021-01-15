package controller

import (
	"hilive/context"
	"hilive/guard"
	"net/http"
)

// NewForm 新增資料POST功能
func (h *Handler) NewForm(ctx *context.Context) {
	param := guard.GetNewForm(ctx)

	err := param.Panel.InsertData(param.MultiForm.Value)
	if err != nil {
		h.showNewForm(ctx, err.Error(), param.Panel, param.Param.GetRouteParamStr(), param.Prefix)
		ctx.AddHeader("X-PJAX-Url", param.Path+"/new"+param.Param.GetRouteParamStr())
		return
	}

	buf := h.showTable(ctx, param.Param, param.Panel, param.Prefix, "")
	ctx.HTML(http.StatusOK, buf.String())
	ctx.AddHeader("X-PJAX-Url", param.Path+param.Param.GetRouteParamStr())
}
