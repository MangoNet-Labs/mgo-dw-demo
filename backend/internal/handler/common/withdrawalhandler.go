package common

import (
	"net/http"
	"user/common/response"
	"user/middleware"

	"github.com/zeromicro/go-zero/rest/httpx"
	"user/internal/logic/common"
	"user/internal/svc"
	"user/internal/types"
)

func WithdrawalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WithdrawalReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailJson(w, err.Error(), 7)
			return
		}

		authUser, ok := middleware.GetAuthUser(r.Context())
		if !ok {
			response.FailJson(w, "unauthorized", 7)
			return
		}

		lBalance := common.NewBalanceLogic(r.Context(), svcCtx)
		respBalance, err := lBalance.Balance(&types.GetBalanceByCoinReq{
			ChainName: req.ChainName,
		}, authUser)
		if err != nil {
			response.FailJson(w, err.Error(), 7)
		}
		l := common.NewWithdrawalLogic(r.Context(), svcCtx)

		var resp *types.WithdrawalResp
		switch req.ChainName {
		case "mgo":
			resp, err = l.Withdrawal(&req, authUser, respBalance.TotalBalance)
		case "sol":
			resp, err = l.WithdrawalSol(&req, authUser, respBalance.TotalBalance)
		default:
			response.FailJson(w, "Parameter error", 7)
			return
		}
		if err != nil {
			response.FailJson(w, err.Error(), 7)
		} else {
			response.OkJson(w, resp)
		}
	}
}
