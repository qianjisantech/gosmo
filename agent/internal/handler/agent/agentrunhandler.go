package agent

import (
	"net/http"

	"agent/internal/logic/agent"
	"agent/internal/svc"
	"agent/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AgentRunHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AgentRunRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := agent.NewAgentRunLogic(r.Context(), svcCtx)
		resp, err := l.AgentRun(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
