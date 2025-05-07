package traffic

import (
	"net/http"

	"github.com/qianjisantech/gosmo/internal/logic/traffic"
	"github.com/qianjisantech/gosmo/internal/svc"
	"github.com/qianjisantech/gosmo/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TrafficPoolQueryPageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TrafficPoolQueryPageRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := traffic.NewTrafficPoolQueryPageLogic(r.Context(), svcCtx)
		resp, err := l.TrafficPoolQueryPage(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
