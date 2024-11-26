package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/rpcox/grpc-protobuf/pkg/job"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	branch   string
	commit   string
	tool     string = `controller`
	version  string = `0.1.0`
	_addr           = flag.String("addr", "localhost:10101", "the address to connect to")
	_jobList        = flag.String("job-list", "job-list", "Identify the job list")
	_log            = flag.String("log", "local.log", tool+" log file")
	_timeOut        = flag.Int("timeout", 5, "context time-out")
	_version        = flag.Bool("version", false, "Display version and exit")
	Debug           = flag.Bool("debug", false, "Enable debug logging")
)

/*func FindServiceAddr(service []string) {

}*/

func NewJob(jobType string) (*job.JobResponse, error) {
	conn, err := grpc.NewClient(*_addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close()

	jrc := job.NewOrderClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*_timeOut)*time.Second)
	defer cancel()

	order := job.JobRequest{Id: job.Order(), JobType: jobType, Device: "mydevice", Issued: time.Now().Unix()}
	resp, err := jrc.Send(ctx, &order)

	return resp, err
}

func Initialize() {
	flag.Parse()
	Version(*_version)
	StartLogging(*_log, nil)
	jobs, err := LoadJobs(*_jobList)
	if err != nil {
		log.Fatal("LoadJobs():", err)
	}
	log.Printf("loading %d jobs\n", len(*jobs))
	tm := TickerMap(jobs)
	log.Printf("loading %d tickers\n", len(tm))
	log.Println("intervals - ", tm)
}

func main() {
	Initialize()
	resp, err := NewJob("state")
	if err != nil {
		// don't do fatal
		log.Fatalf("%v: %v", resp.String(), err)
	}

	//log.Printf("device: %s\nissued: %d\nstart: %d\nend: %d\nduration: %d\n",
	//	resp.GetDevice(), resp.GetIssued(), resp.GetStart(), resp.GetEnd(), resp.GetEnd() - resp.GetStart())
	log.Println("completed:", resp.String())
}
