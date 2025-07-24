//go:build pro

package gosmo

import (
	_ "bufio"
	"fmt"
	_ "io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var _ PluginWriter = (*S3Output)(nil)

// S3Output output plugin
type S3Output struct {
	pathTemplate string

	buffer  *FileOutput
	session *session.Session
	config  *FileOutputConfig
	closeC  chan struct{}
}

// NewS3Output constructor for FileOutput, accepts path
func NewS3Output(pathTemplate string, config *FileOutputConfig) *S3Output {
	o := new(S3Output)
	o.pathTemplate = pathTemplate
	o.config = config
	o.config.onClose = o.onBufferUpdate

	if config.BufferPath == "" {
		config.BufferPath = "/tmp"
	}

	rnd := rand.Int63()
	buffer_name := fmt.Sprintf("gor_output_s3_%d_buf_", rnd)

	pathParts := strings.Split(pathTemplate, "/")
	buffer_name += pathParts[len(pathParts)-1]

	if strings.HasSuffix(o.pathTemplate, ".gz") {
		buffer_name += ".gz"
	}

	buffer_path := filepath.Join(config.BufferPath, buffer_name)

	o.buffer = NewFileOutput(buffer_path, config)
	o.connect()

	return o
}

func (o *S3Output) connect() {
	if o.session == nil {
		o.session = session.Must(session.NewSession(awsConfig()))
		log.Println("[S3 Output] S3 connection succesfully initialized")
	}
}

func (o *S3Output) PluginWrite(msg *Message) (n int, err error) {
	return o.buffer.PluginWrite(msg)
}

func (o *S3Output) String() string {
	return "S3 output: " + o.pathTemplate
}

func (o *S3Output) Close() error {
	return o.buffer.Close()
}

func (o *S3Output) keyPath(idx int) (bucket, key string) {
	bucket, key = parseS3Url(o.pathTemplate)

	for name, fn := range dateFileNameFuncs {
		key = strings.Replace(key, name, fn(o.buffer), -1)
	}

	key = setFileIndex(key, idx)

	return
}

func (o *S3Output) onBufferUpdate(path string) {
	svc := s3.New(o.session)
	idx := getFileIndex(path)
	bucket, key := o.keyPath(idx)

	file, _ := os.Open(path)
	// reader := bufio.NewReader(file)

	_, err := svc.PutObject(&s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Printf("[S3 Output] Failed to upload data to %s/%s, %s\n", bucket, key, err)
		os.Remove(path)
		return
	}

	os.Remove(path)

	if o.closeC != nil {
		o.closeC <- struct{}{}
	}
}
