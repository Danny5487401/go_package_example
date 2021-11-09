// Binary client is an example client.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	pb "go_grpc_example/08_grpc/01_grpc_helloworld/proto"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var addr = flag.String("addr", "localhost:50052", "the address to connect to")

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		if e := conn.Close(); e != nil {
			log.Printf("failed to close connection: %s", e)
		}
	}()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "world"})
	if err != nil {
		// 比FromError更加友好
		s := status.Convert(err)
		for _, d := range s.Details() {
			switch info := d.(type) {
			case *epb.QuotaFailure:
				log.Printf("Quota failure: %s", info)
			default:
				log.Printf("Unexpected type: %s", info)
			}
		}
		os.Exit(1)
	}
	log.Printf("Greeting: %s", r.Message)
}
