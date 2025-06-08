package common

import (
	"errors"
	"net/http"
	"user/common/response"
	"user/middleware"

	"user/internal/logic/common"
	"user/internal/svc"
	"user/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func TransactionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TransactionListReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailJson(w, err.Error(), 7)
			return
		}

		authUser, ok := middleware.GetAuthUser(r.Context())
		if !ok {
			httpx.Error(w, errors.New("unauthorized"))
			return
		}

		l := common.NewTransactionLogic(r.Context(), svcCtx)

		var (
			resp *types.TransactionListResp
			err  error
		)
		switch req.ChainName {
		case "mgo":
			resp, err = l.Transaction(&req, authUser)
		case "sol":
			resp, err = l.TransactionSol(&req, authUser)
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
