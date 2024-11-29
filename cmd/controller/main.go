package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rpcox/grpc-protobuf/pkg/job"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const baseTimeout = time.Second

type JobStream struct {
	RequestStream  chan *job.JobRequest
	ResponseStream chan *job.JobResponse
}

var (
	branch      string
	commit      string
	tool        string = `controller`
	version     string = `0.1.0`
	_jobList           = flag.String("job-list", "job-list", "Identify the job list")
	_log               = flag.String("log", "local.log", tool+" log file")
	_svc_report        = flag.String("svc-report", "localhost:10102", "the address to connect to")
	_svc_state         = flag.String("svc-state", "localhost:10101", "the address to connect to")
	_timeOut           = flag.Int("timeout", 5, "context time-out")
	_version           = flag.Bool("version", false, "Display version and exit")
	Debug              = flag.Bool("debug", false, "Enable debug logging")
)

/*func FindServiceAddr(service []string) {

}*/

func SigHandler(sig chan os.Signal, done chan interface{}) {
	for {
		signal := <-sig
		if signal == syscall.SIGINT || signal == syscall.SIGTERM {
			close(done)
			break
		}

	}

	log.Println("exit sig handler")
}

func Initialize() (*[]Job, *[]Job) {
	flag.Parse()
	Version(*_version)
	StartLogging(*_log, nil)
	state, report, err := LoadJobs(*_jobList)
	if err != nil {
		log.Fatal("LoadJobs():", err)
	}

	return state, report
}

func RunJob(order job.JobRequest) {
	var addr string
	switch order.JobType {
	case "state":
		addr = *_svc_state
	case "report":
		addr = *_svc_report
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error: id=%d, device=%s: %v\n", order.Id, order.Device, err)
		return
	}
	defer conn.Close()

	jrc := job.NewOrderClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*_timeOut)*baseTimeout)
	defer cancel()

	resp, err := jrc.Send(ctx, &order)
	if err != nil {
		log.Printf("error: id=%d, device=%s: %v\n", order.Id, order.Device, err)
	} else {
		log.Printf("success: %s\n", resp.String())
	}
}

func JobTicker(interval int, jobs *[]Job, done chan interface{}, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	log.Printf("starting JobTicker(): interval: %d, jobs: %d\n", interval, len(*jobs))

	for {
		select {
		case <-ticker.C:
			for _, j := range *jobs {
				order := job.JobRequest{Id: job.Order(), JobType: j.Type, Device: j.Device, Issued: time.Now().Unix()}
				log.Println("submitting", order)
				go RunJob(order)

			}
		case _, ok := <-done:
			if !ok {
				log.Printf("exiting JobTicker(): interval: %d, jobs: %d\n", interval, len(*jobs))
				wg.Done()
				return
			}
		}
	}

}

func main() {
	var wg sync.WaitGroup
	done := make(chan interface{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go SigHandler(sig, done)

	stateJobs, reportJobs := Initialize()

	stm := StateTickerMap(stateJobs)
	for interval, _ := range stm {
		jobs := JobsByInterval(interval, stateJobs)
		wg.Add(1)
		go JobTicker(interval, jobs, done, &wg)
	}
	stateJobs = nil // done w/ these structures
	stm = nil

	wg.Add(1)
	go JobTicker(5, reportJobs, done, &wg)

	wg.Wait()
	log.Println("clean exit")
}
