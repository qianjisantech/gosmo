//go:build !pro

package gosmo

import (
	"errors"
	"time"

	"github.com/qianjisantech/gosmo/internal/size"
)

var _ PluginWriter = (*BinaryOutput)(nil)

// BinaryOutputConfig struct for holding binary output configuration
type BinaryOutputConfig struct {
	Workers        int           `json:"output-binary-workers"`
	Timeout        time.Duration `json:"output-binary-timeout"`
	BufferSize     size.Size     `json:"output-tcp-response-buffer"`
	Debug          bool          `json:"output-binary-debug"`
	TrackResponses bool          `json:"output-binary-track-response"`
}

// BinaryOutput plugin manage pool of workers which send request to replayed server
// By default workers pool is dynamic and starts with 10 workers
// You can specify fixed number of workers using `--output-tcp-workers`
type BinaryOutput struct {
	address string
}

// NewBinaryOutput constructor for BinaryOutput
// Initialize workers
func NewBinaryOutput(address string, config *BinaryOutputConfig) PluginReadWriter {
	return &BinaryOutput{address: address}
}

// PluginWrite writes a message to this plugin
func (o *BinaryOutput) PluginWrite(msg *Message) (n int, err error) {
	return 0, errors.New("binary output is only available in PRO version")
}

// PluginRead reads a message from this plugin
func (o *BinaryOutput) PluginRead() (*Message, error) {
	return nil, errors.New("binary output is only available in PRO version")
}

func (o *BinaryOutput) String() string {
	return "Binary output: " + o.address + " (PRO version required)"
}

// Close closes this plugin for reading
func (o *BinaryOutput) Close() error {
	return nil
}
