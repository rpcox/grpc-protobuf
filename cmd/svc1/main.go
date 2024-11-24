package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"github.com/rpcox/grpc-protobuf/pkg/job"
)

var (
        port = flag.Int("port", 10101, "the port to listen to")
)

type server struct {
	job.UnimplementedOrderServer
}

func (s *server) Send(_ context.Context, in *job.JobRequest) (*job.JobResponse, error) {
	log.Printf("Received: %v", in.GetId())
	r := job.JobResponse{}
	r.Id = in.GetId()
	r.JobType = in.GetJobType()
	r.Device = in.GetDevice()
	r.Issued = in.GetIssued()
	t0 := time.Now().Unix()
	r.Start  = t0
        r.End = t0 + 1000

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
