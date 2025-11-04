package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"x-operation/internal/logic"
	"x-operation/internal/svc"
	"x-operation/internal/types"
)

func X_operationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewX_operationLogic(r.Context(), svcCtx)
		resp, err := l.X_operation(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
