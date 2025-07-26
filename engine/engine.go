// Gor is simple http traffic replication tool written in Go. Its main goal to replay traffic from production servers to staging and dev environments.
// Now you can test your code on real user sessions in an automated and repeatable fashion.
package engine

import (
	"flag"
	"fmt"
	recorder "github.com/qianjisantech/gosmo"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write memory profile to this file")
)

type InputRAWConfig struct {
	TrackResponse bool
}

type RecorderEngine struct {
	Port string

	OutputStdout      bool
	InputRAWConfig    InputRAWConfig
	OutputKafkaConfig recorder.OutputKafkaConfig
	// 新增字段
	closeCh      chan struct{}
	emitter      *recorder.Emitter
	plugins      *recorder.InOutPlugins
	stopOnce     sync.Once
	isRunning    bool
	runningMutex sync.Mutex
}

func (rg *RecorderEngine) Start() error {
	rg.runningMutex.Lock()
	defer rg.runningMutex.Unlock()

	if rg.isRunning {
		return fmt.Errorf("engine is already running")
	}

	rg.closeCh = make(chan struct{})
	recorder.Settings.InputRAW = []string{rg.Port}
	recorder.Settings.OutputStdout = rg.OutputStdout
	recorder.Settings.InputRAWConfig.TrackResponse = rg.InputRAWConfig.TrackResponse
	recorder.Settings.OutputKafkaConfig = rg.OutputKafkaConfig
	flag.Parse()
	recorder.CheckSettings()
	rg.plugins = recorder.NewPlugins()

	log.Printf("[PPID %d and PID %d] Version:%s\n", os.Getppid(), os.Getpid(), recorder.VERSION)
	log.Printf("录制%v", rg.plugins.Inputs)
	log.Printf("输出%v", rg.plugins.Outputs)

	if len(rg.plugins.Inputs) == 0 || len(rg.plugins.Outputs) == 0 {
		return fmt.Errorf("required at least 1 input and 1 output")
	}

	// 性能分析
	if *memprofile != "" {
		profileMEM(*memprofile)
	}
	if *cpuprofile != "" {
		profileCPU(*cpuprofile)
	}

	// pprof 服务
	if recorder.Settings.Pprof != "" {
		go func() {
			log.Println(http.ListenAndServe(recorder.Settings.Pprof, nil))
		}()
	}

	// 启动 emitter
	rg.emitter = recorder.NewEmitter()
	go rg.emitter.Start(rg.plugins, recorder.Settings.Middleware)

	// 设置超时自动关闭
	if recorder.Settings.ExitAfter > 0 {
		log.Printf("Running gor for a duration of %s\n", recorder.Settings.ExitAfter)
		time.AfterFunc(recorder.Settings.ExitAfter, func() {
			log.Printf("gor run timeout %s\n", recorder.Settings.ExitAfter)
			rg.Stop()
		})
	}

	rg.isRunning = true
	return nil
}

// Stop 停止录制引擎
func (rg *RecorderEngine) Stop() error {
	rg.runningMutex.Lock()
	defer rg.runningMutex.Unlock()

	if !rg.isRunning {
		return nil
	}

	var stopErr error
	rg.stopOnce.Do(func() {
		log.Println("Stopping recorder engine...")

		// 关闭 emitter
		if rg.emitter != nil {
			rg.emitter.Close()
		}

		// 关闭 plugins
		if rg.plugins != nil {
			for _, input := range rg.plugins.Inputs {
				if closer, ok := input.(io.Closer); ok {
					if err := closer.Close(); err != nil {
						stopErr = fmt.Errorf("input close error: %v", err)
					}
				}
			}
			for _, output := range rg.plugins.Outputs {
				if closer, ok := output.(io.Closer); ok {
					if err := closer.Close(); err != nil {
						stopErr = fmt.Errorf("output close error: %v", err)
					}
				}
			}
		}

		// 关闭通道
		if rg.closeCh != nil {
			close(rg.closeCh)
		}

		rg.isRunning = false
		log.Println("Recorder engine stopped")
	})

	return stopErr
}

// IsRunning 检查引擎是否在运行
func (rg *RecorderEngine) IsRunning() bool {
	rg.runningMutex.Lock()
	defer rg.runningMutex.Unlock()
	return rg.isRunning
}

func profileCPU(cpuprofile string) {
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)

		time.AfterFunc(30*time.Second, func() {
			pprof.StopCPUProfile()
			f.Close()
		})
	}
}

func profileMEM(memprofile string) {
	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Fatal(err)
		}
		time.AfterFunc(30*time.Second, func() {
			pprof.WriteHeapProfile(f)
			f.Close()
		})
	}
}
