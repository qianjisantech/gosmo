package svc

import (
	_ "github.com/IBM/sarama"
	"github.com/qianjisantech/gosmo/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config: c,
	}
}
