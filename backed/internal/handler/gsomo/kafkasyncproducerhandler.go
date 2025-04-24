package gsomo

import (
	"net/http"

	"github.com/qianjisantech/gosmo/internal/logic/gsomo"
	"github.com/qianjisantech/gosmo/internal/svc"
	"github.com/qianjisantech/gosmo/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func KafkaSyncProducerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.KafkaSyncProducerRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := gsomo.NewKafkaSyncProducerLogic(r.Context(), svcCtx)
		resp, err := l.KafkaSyncProducer(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
