package gosmo

import "fmt"

type StdoutOutput struct {
}

func (s *StdoutOutput) PluginWrite(msg *Message) (int, error) {
	if msg == nil {
		return 0, fmt.Errorf("nil message")
	}

	// 打印消息数据到标准输出
	n, err := fmt.Println(string(msg.Data))
	if err != nil {
		return 0, err
	}

	// 返回写入的字节数（减去换行符）
	return n - 1, nil
}
