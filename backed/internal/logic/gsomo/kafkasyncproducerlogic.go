package gsomo

import (
	"context"

	"github.com/qianjisantech/gosmo/internal/svc"
	"github.com/qianjisantech/gosmo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type KafkaSyncProducerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewKafkaSyncProducerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KafkaSyncProducerLogic {
	return &KafkaSyncProducerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *KafkaSyncProducerLogic) KafkaSyncProducer(req *types.KafkaSyncProducerRequest) (resp *types.KafkaSyncProducerResp, err error) {
	// todo: add your logic here and delete this line

	return
}
