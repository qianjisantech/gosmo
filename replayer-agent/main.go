package main

import (
	"context"
	"flag"
	"runtime/debug"
	"strconv"

	"github.com/qianjisantech/gosmo/replayer-agent/common/global"
	"github.com/qianjisantech/gosmo/replayer-agent/common/handlers/conf"
	"github.com/qianjisantech/gosmo/replayer-agent/common/handlers/httpclient"
	"github.com/qianjisantech/gosmo/replayer-agent/common/handlers/httpserv"
	"github.com/qianjisantech/gosmo/replayer-agent/common/handlers/ignore"
	"github.com/qianjisantech/gosmo/replayer-agent/common/handlers/limit"
	"github.com/qianjisantech/gosmo/replayer-agent/common/handlers/module"
	"github.com/qianjisantech/gosmo/replayer-agent/common/handlers/outbound"
	"github.com/qianjisantech/gosmo/replayer-agent/common/handlers/template"
	"github.com/qianjisantech/gosmo/replayer-agent/common/handlers/tlog"
	"github.com/qianjisantech/gosmo/replayer-agent/model/nuwaplt"
	"github.com/qianjisantech/gosmo/replayer-agent/router"
)

func init() {
	flag.BoolVar(&global.FlagHandler.EnableCursor, "cursor", false, "enable cursor for matching alg")
	flag.IntVar(&global.FlagHandler.Parallel, "parallel", 10, "set max parallel num for replaying")
	flag.Parse()

	conf.Init("")
	tlog.Init()
	httpclient.Init()

	ignore.Init()
	limit.Init()
	outbound.Init()
	module.Init()
	nuwaplt.Reload()

	template.Init()
	router.Init()
	httpserv.Init()
}

func main() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				tlog.Handler.Errorf(context.Background(), tlog.DLTagUndefined, "panic in %s goroutine||errmsg=%s||stack info=%s", "sharingan", err, strconv.Quote(string(debug.Stack())))
			}
		}()
	}()

	httpserv.Run()
}
