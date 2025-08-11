package gosmo

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
)

func TestOutputKafkaRAW(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer := mocks.NewAsyncProducer(t, config)
	producer.ExpectInputAndSucceed()

	output := NewKafkaOutput("", &OutputKafkaConfig{
		producer: producer,
		Topic:    "test",
		UseJSON:  false,
	}, nil)

	output.PluginWrite(&Message{Meta: []byte("1 2 3\n"), Data: []byte("GET / HTTP1.1\r\nHeader: 1\r\n\r\n")})

	resp := <-producer.Successes()

	data, _ := resp.Value.Encode()

	if string(data) != "1 2 3\nGET / HTTP1.1\r\nHeader: 1\r\n\r\n" {
		t.Errorf("Message not properly encoded: %q", data)
	}
}

func TestOutputKafkaJSON(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer := mocks.NewAsyncProducer(t, config)
	producer.ExpectInputAndSucceed()

	output := NewKafkaOutput("", &OutputKafkaConfig{
		producer: producer,
		Topic:    "gosmo",
		UseJSON:  true,
	}, nil)

	output.PluginWrite(&Message{Meta: []byte("1 2 3\n"), Data: []byte("GET / HTTP1.1\r\nHeader: 1\r\n\r\n")})

	resp := <-producer.Successes()

	data, _ := resp.Value.Encode()

	if string(data) != `{"Req_URL":"","Req_Type":"1","Req_ID":"2","Req_Ts":"3","Req_Method":"GET"}` {
		t.Error("Message not properly encoded: ", string(data))
	}
}
func TestRealOutputKafkaJSON(t *testing.T) {
	// Kafka 配置
	topic := "polaris"
	host := "10.24.11.21:30320,10.24.11.22:30321,10.24.11.21:30322"

	t.Log("=== 开始 Kafka 生产者测试 ===")
	t.Logf("Brokers: %v", host)
	t.Logf("Topic: %s", topic)

	// 配置 Sarama
	config := sarama.NewConfig()

	// 生产相关配置
	config.Producer.Return.Successes = true          // 需要返回成功确认
	config.Producer.RequiredAcks = sarama.WaitForAll // 等待所有副本确认
	config.Producer.Retry.Max = 5                    // 最大重试次数

	// 幂等生产者配置
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1 // 幂等生产者必须设置为1

	// 网络相关配置
	config.Net.DialTimeout = 15 * time.Second
	config.Net.ReadTimeout = 15 * time.Second
	config.Net.WriteTimeout = 15 * time.Second
	config.Net.KeepAlive = 30 * time.Second

	// 元数据配置
	config.Metadata.Retry.Max = 3
	config.Metadata.Retry.Backoff = 1 * time.Second
	config.Metadata.RefreshFrequency = 10 * time.Minute

	// 客户端标识
	config.ClientID = "gosmo-test-client-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	// 创建 Kafka 生产者
	producer, err := sarama.NewAsyncProducer(strings.Split(host, ","), config)
	if err != nil {
		t.Fatalf("创建 Kafka 生产者失败: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			t.Errorf("关闭生产者时出错: %v", err)
		}
	}()

	// 处理发送结果
	done := make(chan bool)
	go func() {
		defer close(done)
		for {
			select {
			case success := <-producer.Successes():
				t.Logf("消息发送成功，分区：%d，偏移量：%d",
					success.Partition, success.Offset)
				return
			case err := <-producer.Errors():
				t.Errorf("消息发送失败: %v", err)
				return
			}
		}
	}()

	// 创建 Kafka 输出插件
	output := NewKafkaOutput("", &OutputKafkaConfig{
		producer: producer,
		Topic:    topic,
		UseJSON:  true,
		Host:     host,
	}, nil)

	// 测试消息
	testMsg := &Message{
		Meta: []byte("1 2 3\n"),
		Data: []byte("GET / HTTP1.1\r\nHeader: 1\r\n\r\n"),
	}

	// 写入消息
	t.Log("尝试写入测试消息...")
	if _, err := output.PluginWrite(testMsg); err != nil {
		t.Fatalf("写入消息失败: %v", err)
	}

	// 等待发送结果
	select {
	case <-done:
		t.Log("消息发送处理完成")
	case <-time.After(20 * time.Second):
		t.Fatal("等待消息发送确认超时")
	}

	t.Log("=== 测试完成 ===")
}
