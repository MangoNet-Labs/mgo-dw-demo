package common

import (
	"net/http"
	"user/common/response"
	"user/internal/logic/common"
	"user/middleware"

	"github.com/zeromicro/go-zero/rest/httpx"
	"user/internal/svc"
	"user/internal/types"
)

func BalanceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetBalanceByCoinReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailJson(w, err.Error(), 7)
			return
		}
		authUser, ok := middleware.GetAuthUser(r.Context())
		if !ok {
			response.FailJson(w, "unauthorized", 7)
			return
		}
		l := common.NewBalanceLogic(r.Context(), svcCtx)
		resp, err := l.Balance(&req, authUser)
		if err != nil {
			response.FailJson(w, err.Error(), 7)
		} else {
			response.OkJson(w, resp)
		}
	}
}
