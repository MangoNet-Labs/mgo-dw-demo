package common

import (
	"net/http"
	"user/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"user/internal/logic/common"
	"user/internal/svc"
	"user/internal/types"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.FailJson(w, err.Error(), 7)
			return
		}

		l := common.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			response.FailJson(w, err.Error(), 7)
		} else {
			response.OkJson(w, resp)
		}
	}
}
