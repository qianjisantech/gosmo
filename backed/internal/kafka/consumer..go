package kafka

import (
	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConsumerHandler struct {
	MsgChan chan<- *sarama.ConsumerMessage // 只写通道
	Logger  logx.Logger
}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	h.Logger.Info("Consumer handler setup")
	return nil
}

func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	h.Logger.Info("Consumer handler cleanup")
	return nil
}

func (h *ConsumerHandler) ConsumeClaim(
	sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		select {
		case h.MsgChan <- msg: // 发送到服务上下文的消息通道
			sess.MarkMessage(msg, "") // 提交偏移量
		default:
			h.Logger.Error("Message channel full, dropping message")
		}
	}
	return nil
}
