package traffic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/qianjisantech/gosmo/internal/common/errorx"
	"github.com/qianjisantech/gosmo/internal/svc"
	"github.com/qianjisantech/gosmo/internal/types"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type TrafficPoolQueryPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTrafficPoolQueryPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrafficPoolQueryPageLogic {
	return &TrafficPoolQueryPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type ElasticsearchResponse struct {
	Hits struct {
		Hits []struct {
			Index   string   `json:"_index"`
			Type    string   `json:"_type"`
			ID      string   `json:"_id"`
			Score   float64  `json:"_score"`
			Ignored []string `json:"_ignored,omitempty"`
			Source  struct {
				Timestamp  string            `json:"@timestamp"`
				ReqBody    string            `json:"req_body"`
				ReqHeaders map[string]string `json:"req_headers"`
				ReqID      string            `json:"req_id"`
				ReqMethod  string            `json:"req_method"`
				ReqTS      string            `json:"req_ts"`
				ReqURL     string            `json:"req_url"`
				RespData   struct {
					Headers map[string]string `json:"headers"`
					Body    string            `json:"body"`
					Status  string            `json:"status"`
					TS      string            `json:"ts"`
				} `json:"resp_data"`
				Type string `json:"type"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (l *TrafficPoolQueryPageLogic) TrafficPoolQueryPage(req *types.TrafficPoolQueryPageRequest) (resp *types.TrafficPoolQueryPageResp, err error) {
	// Elasticsearch配置
	esURL := "http://101.201.116.86:9200"
	indexName := "gosmo"

	// 构建请求URL
	url := fmt.Sprintf("%s/%s/_search?filter_path=hits.hits", esURL, indexName)

	// 创建HTTP请求
	request, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return nil, errorx.NewDefaultError(err.Error())
	}

	// 设置请求头
	request.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return nil, errorx.NewDefaultError(err.Error())
	}
	defer response.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return nil, errorx.NewDefaultError(err.Error())
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("解析JSON失败: %v\n", err)
		return nil, errorx.NewDefaultError(err.Error())
	}

	// 打印结果
	fmt.Println("查询结果:")
	prettyJSON, _ := json.MarshalIndent(result, "", "  ")
	var elasticsearchResponse ElasticsearchResponse
	err = json.Unmarshal(prettyJSON, &elasticsearchResponse)
	if err != nil {
		log.Fatalf("JSON解析失败: %v", err)
	}
	var trafficPoolBodys []TrafficPoolBody
	if len(elasticsearchResponse.Hits.Hits) > 0 {
		for _, hit := range elasticsearchResponse.Hits.Hits {
			trafficPoolBodyRequest := &TrafficPoolBodyRequest{
				Body:    hit.Source.ReqBody,
				Headers: hit.Source.ReqHeaders,
			}
			trafficPoolBodyResponse := &TrafficPoolBodyResponse{
				Body:    hit.Source.RespData.Body,
				Headers: hit.Source.RespData.Headers,
			}
			trafficPoolBody := &TrafficPoolBody{
				Index:     hit.Index,
				Id:        hit.ID,
				TimeStamp: hit.Source.Timestamp,
				Host:      "localhost",
				Method:    hit.Source.ReqMethod,
				Url:       hit.Source.ReqURL,
				RequestId: hit.Source.ReqID,
				Status:    hit.Source.RespData.Status,
				Request:   *trafficPoolBodyRequest,
				Response:  *trafficPoolBodyResponse,
			}
			tsNano, _ := strconv.ParseInt(hit.Source.RespData.TS, 10, 64)

			tsSec := tsNano / 1e9 // 转换为秒级时间戳
			trafficPoolBody.RT = strconv.FormatInt(tsSec, 10)
			trafficPoolBodys = append(trafficPoolBodys, *trafficPoolBody)
		}
	} else {
		trafficPoolBodys = []TrafficPoolBody{}
	}
	return &types.TrafficPoolQueryPageResp{
		Success: true,
		Message: "success",
		Data:    trafficPoolBodys,
	}, nil
}

type TrafficPoolBodyRequest struct {
	Body    interface{}       `json:"body"`
	Headers map[string]string `json:"headers"`
}
type TrafficPoolBodyResponse struct {
	Body    interface{}       `json:"body"`
	Headers map[string]string `json:"headers"`
}
type TrafficPoolBody struct {
	Index     string                  `json:"index"`
	Id        string                  `json:"id"`
	TimeStamp string                  `json:"timestamp"`
	Host      string                  `json:"host"`
	Method    string                  `json:"method"`
	Url       string                  `json:"url"`
	RequestId string                  `json:"requestId"`
	Status    string                  `json:"status"`
	RT        string                  `json:"rt"`
	Request   TrafficPoolBodyRequest  `json:"request"`
	Response  TrafficPoolBodyResponse `json:"response"`
}
