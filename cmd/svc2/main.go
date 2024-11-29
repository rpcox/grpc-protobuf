package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/rpcox/grpc-protobuf/pkg/job"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var (
	branch  string
	commit  string
	tool    string = `svc2`
	version string = `0.1.0`
	Debug          = flag.Bool("debug", false, "Enable debug logging")
	_log           = flag.String("log", "local.log", tool+" log file")
	_port          = flag.Int("port", 10102, "Identify the port to listen at")
	//_sdisc = flag.String("sdisc", ":80", "Identify the service discovery agent")
	_version = flag.Bool("version", false, "Display version and exit")
)

func Version(b bool) {
	if b {
		if commit != "" {
			// go build -ldflags="-X main.commit=$(git rev-parse --short HEAD) -X main.branch=$(git branch | sed 's/.*\* //')"
			fmt.Printf("%s v%s (commit:%s branch:%s)\n", tool, version, commit, branch)
		} else {
			// go build
			fmt.Printf("%s v%s\n", tool, version)
		}

		os.Exit(0)
	}
}

// Initialize logging
func StartLogging(fileName string, currFile *os.File) *os.File {
	if currFile != nil {
		log.Println("closing log")
		currFile.Close()
	}

	fh, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil
	}

	log.SetOutput(fh)
	if *Debug {
		log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.LUTC | log.Ldate | log.Ldate)
	} else {
		log.SetFlags(log.Lmicroseconds | log.LUTC | log.Ldate | log.Ldate)
	}

	log.Printf("%s v%s\n", tool, version)
	return fh
}

type server struct {
	job.UnimplementedOrderServer
}

func (s *server) Send(ctx context.Context, in *job.JobRequest) (*job.JobResponse, error) {
	log.Printf("report job %v received\n", in.GetId())
	if p, ok := peer.FromContext(ctx); ok {
		log.Println(p.String())
	}
	r := job.JobResponse{Id: in.GetId(), JobType: in.GetJobType(), Device: in.GetDevice(), Issued: in.GetIssued()}
	r.Start = job.TimeStamp()

	// do something
	time.Sleep(5 * time.Second)
	r.End = job.TimeStamp()
	log.Printf("report job %v completed\n", in.GetId())

	return &r, nil
}

func main() {
	flag.Parse()
	Version(*_version)
	StartLogging(*_log, nil)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *_port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer lis.Close()

	s := grpc.NewServer()
	job.RegisterOrderServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

}
