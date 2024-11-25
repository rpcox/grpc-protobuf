package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
        "google.golang.org/grpc/credentials/insecure"
	"github.com/rpcox/grpc-protobuf/pkg/job"
)

var (
        _addr = flag.String("addr", "localhost:10101", "the address to connect to")
	_timeOut = flag.Int("time-out", 5, "context time-out")

)

func main() {
	flag.Parse()

	// Set up a connection to the server.
        conn, err := grpc.NewClient(*_addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
                log.Fatalf("did not connect: %v", err)
        }

        defer conn.Close()

	jrc := job.NewOrderClient(conn)
	// Contact the server and print out its response.
        ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*_timeOut) * time.Second)
        defer cancel()

	order := job.JobRequest{Id: job.Order(), JobType: "state", Device: "mydevice", Issued: time.Now().Unix()}
        resp, err := jrc.Send(ctx, &order)
	if err != nil {
                log.Fatalf("%v: %v", order.String(), err)
        }

	//log.Printf("device: %s\nissued: %d\nstart: %d\nend: %d\nduration: %d\n",
	//	resp.GetDevice(), resp.GetIssued(), resp.GetStart(), resp.GetEnd(), resp.GetEnd() - resp.GetStart())
	log.Println("as string:", resp.String())
}
