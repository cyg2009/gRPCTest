package main

import (
	"log"
	"os"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "serverlessgo/greet"
)

const (
	address     = "localhost:50051"
	defaultName = `{"name": "into the new world"}`
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	r, err := c.Greet(context.Background(), &pb.GreetRequest{Message: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Respoonse: %s", r.Message)
}
