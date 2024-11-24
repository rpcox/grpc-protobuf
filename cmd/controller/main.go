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
        addr = flag.String("addr", "localhost:10101", "the address to connect to")
)

func main() {
	flag.Parse()

	// Set up a connection to the server.
        conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
                log.Fatalf("did not connect: %v", err)
        }

        defer conn.Close()

	jrc := job.NewOrderClient(conn)
	// Contact the server and print out its response.
        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()

	order := job.JobRequest{Id: 100, JobType: "state", Device: "mydevice", Issued: 1000}
        resp, err := jrc.Send(ctx, &order)
	if err != nil {
                log.Fatalf("could not kick job: %v", err)
        }

	log.Printf("device: %s\nissued: %d\nstart: %d\nend: %d\nduration: %d\n",
		resp.GetDevice(), resp.GetIssued(), resp.GetStart(), resp.GetEnd(), resp.GetEnd() - resp.GetStart())
	log.Println("as string:", resp.String())
}
