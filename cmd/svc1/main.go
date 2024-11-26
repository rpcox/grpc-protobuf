package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"github.com/rpcox/grpc-protobuf/pkg/job"
)

var (
        port = flag.Int("port", 10101, "the port to listen to")
)

type server struct {
	job.UnimplementedOrderServer
}

func (s *server) Send(ctx context.Context, in *job.JobRequest) (*job.JobResponse, error) {
	log.Printf("job received: %v", in.GetId())
	if p, ok := peer.FromContext(ctx); ok {
		//log.Println(p.Addr)
		log.Println(p.String())
	}
	r := job.JobResponse{Id: in.GetId(), JobType: in.GetJobType(), Device: in.GetDevice(), Issued: in.GetIssued()}
	r.Start = job.TimeStamp()

	// do something
	time.Sleep(5 * time.Second)
	r.End = job.TimeStamp()
	log.Printf("job completed: %v", in.GetId())

	return &r, nil
}

func main() {
        flag.Parse()
        lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
        if err != nil {
                log.Fatalf("failed to listen: %v", err)
        }

        defer lis.Close()

	s := grpc.NewServer()
        job.RegisterOrderServer(s, &server{})
        log.Printf("server listening at %v", lis.Addr())
        if err := s.Serve(lis); err != nil {
                log.Fatalf("failed to serve: %v", err)
        }

}
