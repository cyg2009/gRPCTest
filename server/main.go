package main

import (
	"log"
	"net"
	"encoding/json"
	"plugin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "serverlessgo/greet"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

func invoke(target string, evtJSON string) (string, error) {

	 p, err := plugin.Open(target)
        if err != nil {
            return "", err
        }

        greetSymbol, err := p.Lookup("Handler")
        if err != nil {
            return "", err
        }

		trigger := greetSymbol.(func(string) string)
		
		return trigger(evtJSON), nil
}

// Implement greet.GreeterServer.
type server struct{}

// Implements greet.GreeterServer
func (s *server) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetReply, error) {
	m := map[string]string{}
	err := json.Unmarshal([]byte(in.Message), &m)
	if err != nil {	
		return nil, err
	}

	action, ok  := m["action"]
	if ok {
		log.Printf("action: %s", action)

		switch action {
			case "health":
			   return &pb.GreetReply{Message: "OK"}, nil	
			case "invoke": {
			   name, ok := m["function"]
               if ok == false {
				   return &pb.GreetReply{Message: "invoke unknown function"}, nil	
			   }
				
			   evt, ok := m["event"]
			   if ok == false {
				   return &pb.GreetReply{Message: "invoke function with no event"}, nil	
			   }

			   ret, err := invoke("functionplugin.so", evt)
				
			   if err != nil {
				   return &pb.GreetReply{Message: "invoke function with error"}, err	
			   }
			   //we need use a plugin here to invoke it
			   return &pb.GreetReply{Message: "function invoked:" + name + " --- " + ret}, nil	
			}
			   		
			default:
	           return &pb.GreetReply{Message: "Unknown action " + action}, nil			
		}
	}

	name, ok  := m["name"]
	if ok {
		log.Printf("got your name: %s", name)
		return &pb.GreetReply{Message: "Hello " + name}, nil
	}
	return &pb.GreetReply{Message: "Unable to process this message: " + in.Message}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
