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
	engines    = make(map[int]*RecorderEngine)
	enginesMux sync.Mutex
)

type InputRAWConfig struct {
	TrackResponse bool
}

type RecorderEngine struct {
	Port              string
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
	pid          int
}

func (rg *RecorderEngine) Start() (int, error) {
	rg.runningMutex.Lock()
	defer rg.runningMutex.Unlock()

	if rg.isRunning {
		return 0, fmt.Errorf("engine is already running")
	}

	rg.closeCh = make(chan struct{})
	recorder.Settings.InputRAW = []string{rg.Port}
	recorder.Settings.OutputStdout = rg.OutputStdout
	recorder.Settings.InputRAWConfig.TrackResponse = rg.InputRAWConfig.TrackResponse
	recorder.Settings.OutputKafkaConfig = rg.OutputKafkaConfig
	flag.Parse()
	recorder.CheckSettings()
	rg.plugins = recorder.NewPlugins()

	rg.pid = os.Getpid()
	enginesMux.Lock()
	engines[rg.pid] = rg
	enginesMux.Unlock()

	log.Printf("[PPID %d and PID %d] Version:%s\n", os.Getppid(), rg.pid, recorder.VERSION)
	log.Printf("录制%v", rg.plugins.Inputs)
	log.Printf("输出%v", rg.plugins.Outputs)

	if len(rg.plugins.Inputs) == 0 || len(rg.plugins.Outputs) == 0 {
		return 0, fmt.Errorf("required at least 1 input and 1 output")
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

	rg.isRunning = true
	return rg.pid, nil
}

// Stop 停止录制引擎
func (rg *RecorderEngine) Stop(pid int) error {
	enginesMux.Lock()
	rg, exists := engines[pid]
	enginesMux.Unlock()

	if !exists {
		return fmt.Errorf("engine with PID %d not found", pid)
	}

	rg.runningMutex.Lock()
	defer rg.runningMutex.Unlock()

	if !rg.isRunning {
		return nil
	}

	var stopErr error
	rg.stopOnce.Do(func() {
		log.Printf("Stopping recorder engine with PID %d...", pid)

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
		enginesMux.Lock()
		delete(engines, pid)
		enginesMux.Unlock()
		log.Printf("Recorder engine with PID %d stopped", pid)
	})

	return stopErr
}

// IsRunning 检查引擎是否在运行
func (rg *RecorderEngine) IsRunning(pid int) bool {
	enginesMux.Lock()
	rg, exists := engines[pid]
	enginesMux.Unlock()

	if !exists {
		return false
	}

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
			err := pprof.WriteHeapProfile(f)
			if err != nil {
				return
			}
			err = f.Close()
			if err != nil {
				return
			}
		})
	}
}
