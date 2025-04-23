package search

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/didi/sharingan/replayer-agent/common/handlers/conf"
	"github.com/didi/sharingan/replayer-agent/common/handlers/httpclient"
	"github.com/didi/sharingan/replayer-agent/common/handlers/tlog"
	"github.com/didi/sharingan/replayer-agent/idl"
	"github.com/didi/sharingan/replayer-agent/model/esmodel"
	"github.com/didi/sharingan/replayer-agent/model/nuwaplt"
	"github.com/didi/sharingan/replayer-agent/utils/helper"

	jsoniter "github.com/json-iterator/go"
)

type SearchRecords []*SearchRecord

//Len()
func (s SearchRecords) Len() int {
	return len(s)
}

//Less():成绩将有低到高排序
func (s SearchRecords) Less(i, j int) bool {
	return s[i].Timestamp > s[j].Timestamp
}

//Swap()
func (s SearchRecords) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

/**
 * 前端查询ES数据接口(各种查询条件)，超时10s
 * @Return 多条session，封装成前端展示数据结构
 */
func Search(ctx context.Context, req *idl.SearchReq) []*SearchRecord {
	req = SearchPreHandle(req)

	// read from es
	if conf.Handler.GetString("es_url.default") != "" {
		return SearchFromES(ctx, req)
	}

	// read from local
	stat, contents := GetTrafficFromLocal(ctx, req)
	if stat != 0 {
		return nil
	}

	return SearchHandleLocal(contents, req)
}

// SearchPreHandle 处理Search接口参数
func SearchPreHandle(req *idl.SearchReq) *idl.SearchReq {
	if len(req.Date) == 1 {
		req.Start = req.Date[0]
		req.End = time.Now().Format("2006-01-02")
	} else if len(req.Date) == 2 {
		req.Start = req.Date[0]
		req.End = req.Date[1]
	}

	return req
}

// SearchFromES 从es读取数据，并处理按Search接口返回
func SearchFromES(ctx context.Context, req *idl.SearchReq) []*SearchRecord {
	flowList := make([]*SearchRecord, 0)
	body, qErr := Query(ctx, req, 0)
	if qErr != nil {
		return flowList
	}

	if queryForID(req.Field) {
		return retrieveSessionIds(ctx, body, flowList)
	}

	return retrieveSessions(ctx, req.Project, body, flowList)
}

// SearchHandleLocal 处理本地流量符合Search接口
func SearchHandleLocal(contents []string, req *idl.SearchReq) []*SearchRecord {
	flowList := make([]*SearchRecord, 0)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//原始流量格式
	for _, flow := range contents {
		traffic := &esmodel.Session{}
		err := json.Unmarshal([]byte(flow), traffic)
		if err != nil {
			continue
		}

		flowList = AppendFlowList(flowList, req, *traffic)
		sort.Sort(SearchRecords(flowList))
	}

	return flowList
}

// GetTrafficFromLocal 从本地读取流量
func GetTrafficFromLocal(ctx context.Context, req *idl.SearchReq) (int, []string) {
	//读取配置文件conf/traffic/{project}
	contents, err := helper.ReadLines(conf.Root + "/conf/traffic/" + req.Project)
	if err != nil {
		tlog.Handler.Errorf(ctx, tlog.DLTagUndefined, "Failed to read /conf/traffic/"+req.Project+", err="+err.Error())
		return 1, nil
	}
	return 0, contents
}

//AppendFlowList 筛选流量并格式化，并追加到flowList数组
func AppendFlowList(flowList []*SearchRecord, req *idl.SearchReq, session esmodel.Session) []*SearchRecord {
	if valide, data := filterTraffic(req, session); valide {
		if queryForID(req.Field) {
			flowList = append(flowList, &SearchRecord{SessionId: data.SessionId})
		} else {
			flowList = append(flowList, &SearchRecord{Project: data.Project, SessionId: data.SessionId, Timestamp: data.Timestamp, Req: data.Req, Res: data.Res, Actions: data.Actions})
		}
	}
	return flowList
}

//filterTraffic 过滤配置文件conf/traffic/{project}下的流量
func filterTraffic(req *idl.SearchReq, session esmodel.Session) (bool, *SearchRecord) {
	traffic := handleOneSessionRaw(context.TODO(), session, req.Project)
	if traffic == nil {
		return false, nil
	}

	// 线下回放不过滤时间
	// t := strings.Split(traffic.Timestamp, "T")
	// if req.Start > t[0] || req.End < t[0] {
	// 	return false, nil
	// }

	if req.InboundRequest != "" {
		if !strings.Contains(traffic.Req, req.InboundRequest) {
			return false, nil
		}
	}

	if req.InboundResponse != "" {
		if !strings.Contains(traffic.Res, req.InboundResponse) {
			return false, nil
		}
	}

	if req.OutboundRequest != "" || req.OutboundResponse != "" || req.Apollo != "" {
		signReq, signRes, signApollo := true, true, true
		if req.OutboundRequest != "" {
			signReq = false
		}
		if req.OutboundResponse != "" {
			signRes = false
		}
		if req.Apollo != "" {
			signApollo = false
		}
		for _, subReq := range traffic.Actions {
			if req.OutboundRequest != "" && strings.Contains(subReq["req"], req.OutboundRequest) {
				signReq = true
			}
			if req.OutboundResponse != "" && strings.Contains(subReq["res"], req.OutboundResponse) {
				signRes = true
			}
			if req.Apollo != "" && strings.Contains(subReq["apollo"], req.Apollo) {
				signApollo = true
			}
			if signReq && signRes && signApollo {
				break
			}
		}

		if !signRes || !signReq || !signApollo {
			return false, nil
		}
	}

	if req.SessionId != "" {
		if !strings.Contains(traffic.SessionId, req.SessionId) {
			return false, nil
		}
	}

	return true, traffic
}

/**
 * 回放某个具体session时，后端查询ES数据接口(根据SessionId查询)，超时3s
 * @Return 一个session，原始数据格式
 */
func GetRawSessions(ctx context.Context, req *idl.SearchReq) *esmodel.Session {
	//优先读取es地址
	if conf.Handler.GetString("es_url.default") != "" {
		return GetRawSessionFromES(ctx, req)
	}

	// read from local
	stat, contents := GetTrafficFromLocal(ctx, req)
	if stat != 0 {
		return nil
	}

	return GetRawSessionHandleLocal(ctx, contents, req)
}

func GetRawSessionFromES(ctx context.Context, req *idl.SearchReq) *esmodel.Session {
	body, err := Query(ctx, req, 3*time.Second)
	if err != nil {
		tlog.Handler.Errorf(ctx, tlog.DLTagUndefined, "errmsg= GetRawSessions failed from es ||err=%s", err.Error())
		return nil
	}
	sessions, err := esmodel.RetrieveSessions(body)
	if err != nil {
		tlog.Handler.Errorf(ctx, tlog.DLTagUndefined, "errmsg= GetRawSessions failed at retrieve session from es ||err=%s", err.Error())
		return nil
	} else if len(sessions) == 0 {
		tlog.Handler.Errorf(ctx, tlog.DLTagUndefined, "errmsg= GetRawSessions got empty session from es!")
		return nil
	}

	return &sessions[0]
}

// GetRawSessionHandleLocal 处理本地流量符合GetRawSessions接口
func GetRawSessionHandleLocal(ctx context.Context, contents []string, req *idl.SearchReq) *esmodel.Session {
	session := new(esmodel.Session)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for _, flow := range contents {
		traffic := &esmodel.Session{}
		err := json.Unmarshal([]byte(flow), traffic)
		if err != nil {
			tlog.Handler.Warnf(ctx, tlog.DLTagUndefined, "errmsg= GetRawSessions failed at unmarshal from conf with origin struct ||err=%s", err.Error())
			continue
		}
		if traffic.SessionId == req.SessionId {
			session = traffic
			break
		}
	}

	if session == nil {
		tlog.Handler.Warnf(ctx, tlog.DLTagUndefined, "errmsg=receive zero session from es or conf")
	}

	return session
}

func Query(ctx context.Context, req *idl.SearchReq, timeout time.Duration) ([]byte, error) {
	dsl, err := toDSL(req)
	if err != nil {
		tlog.Handler.Errorf(ctx, tlog.DLTagUndefined, "errmsg=generate dsl failed||err=%s", err)
		return nil, err
	}

	return doQuery(ctx, dsl, timeout)
}

func doQuery(ctx context.Context, data []byte, timeout time.Duration) (body []byte, err error) {
	headers := map[string]string{
		"Content-Type": "application/json",
		"kbn-xsrf":     "1",
	}

	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	for i := 0; i < 3; i++ {
		_, body, err = httpclient.Handler.Post(ctx, resolveURL(ctx), data, timeout, headers)
		if err == nil && len(body) > 180 {
			return
		}
	}
	return
}

func resolveURL(ctx context.Context) string {
	url := conf.Handler.GetString("es_url.default")
	depart, ok := ctx.Value(nuwaplt.KDepartment).(string)
	if ok && depart != nuwaplt.DefaultDepartment {
		otherUrl := conf.Handler.GetString("es_url." + depart)
		if otherUrl != "" {
			url = otherUrl
		}
	}
	return url
}
