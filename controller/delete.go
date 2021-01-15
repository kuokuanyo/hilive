package controller

import (
	"hilive/context"
	"hilive/guard"
	"hilive/modules/auth"
	"hilive/modules/response"
)

// Delete 刪除資料POST功能
func (h *Handler) Delete(ctx *context.Context) {
	param := guard.GetDeleteParam(ctx)

	err := h.GetTable(ctx, param.Prefix).DeleteData(param.ID)
	if err != nil {
		response.Error(ctx, "delete fail")
		panic(err)
	}

	response.OkWithData(ctx, map[string]interface{}{
		"token": auth.ConvertInterfaceToTokenService(h.services.Get("token_csrf_helper")).AddToken(),
	})
}
