//go:build !pro

package gosmo

import (
	"errors"
	"fmt"
)

// S3Output output plugin
type S3Output struct{}

// NewS3Output constructor for FileOutput, accepts path
func NewS3Output(pathTemplate string, config *FileOutputConfig) *S3Output {
	fmt.Println("S3 output is only available in the pro version")
	return &S3Output{}
}

func (o *S3Output) PluginWrite(msg *Message) (n int, err error) {
	return 0, errors.New("S3 output is only available in the pro version")
}

func (o *S3Output) String() string {
	return "S3 output (pro version only)"
}

func (o *S3Output) Close() error {
	return errors.New("S3 output is only available in the pro version")
}
