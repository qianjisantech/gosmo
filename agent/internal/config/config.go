package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	ElasticSearch struct {
		Hosts    string
		Username string
		Password string
		Index    string
	}
}
