package svc

import (
	"agent/internal/config"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config   config.Config
	ESClient *elasticsearch.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{c.ElasticSearch.Hosts},
		//Username:  c.Elasticsearch.Username,
		//Password:  c.Elasticsearch.Password,
	})
	if err != nil {
		logx.Must(err)
	}
	return &ServiceContext{
		Config:   c,
		ESClient: esClient,
	}
}
