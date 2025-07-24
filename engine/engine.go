// Gor is simple http traffic replication tool written in Go. Its main goal to replay traffic from production servers to staging and dev environments.
// Now you can test your code on real user sessions in an automated and repeatable fashion.
package engine

import (
	"flag"
	"fmt"
	recorder "github.com/qianjisantech/gosmo"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
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
	Port           string
	OutputStdout   bool
	InputRAWConfig InputRAWConfig
}

func (rg *RecorderEngine) Start() error {
	var plugins *recorder.InOutPlugins
	recorder.Settings.InputRAW = []string{rg.Port}
	recorder.Settings.OutputStdout = rg.OutputStdout
	recorder.Settings.InputRAWConfig.TrackResponse = rg.InputRAWConfig.TrackResponse
	flag.Parse()
	recorder.CheckSettings()
	plugins = recorder.NewPlugins()

	log.Printf("[PPID %d and PID %d] Version:%s\n", os.Getppid(), os.Getpid(), recorder.VERSION)
	log.Printf("录制%v", plugins.Inputs)
	log.Printf("输出%v", plugins.Outputs)
	if len(plugins.Inputs) == 0 || len(plugins.Outputs) == 0 {

		log.Fatal("Required at least 1 input and 1 output")
		return fmt.Errorf("Required at least 1 input and 1 output  %v", plugins)
	}

	if *memprofile != "" {
		profileMEM(*memprofile)
	}

	if *cpuprofile != "" {
		profileCPU(*cpuprofile)
	}

	if recorder.Settings.Pprof != "" {
		go func() {
			log.Println(http.ListenAndServe(recorder.Settings.Pprof, nil))
		}()
	}

	closeCh := make(chan int)
	emitter := recorder.NewEmitter()
	go emitter.Start(plugins, recorder.Settings.Middleware)
	if recorder.Settings.ExitAfter > 0 {
		log.Printf("Running gor for a duration of %s\n", recorder.Settings.ExitAfter)

		time.AfterFunc(recorder.Settings.ExitAfter, func() {
			log.Printf("gor run timeout %s\n", recorder.Settings.ExitAfter)
			close(closeCh)
		})
	}
	return nil
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
