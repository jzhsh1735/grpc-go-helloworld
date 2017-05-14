package main

import (
	"flag"
	"fmt"
	"io"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "myexample.com/helloworld/hello"
)

var (
	count = flag.Int("count", 5, "The server reply count")
	reply = flag.String("reply", "HelloClient!", "The server reply")
	port = flag.Int("port", 12345, "The server port")
)

type helloServiceServer struct {
	count int
	reply string
}

func (s *helloServiceServer) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Reply: s.reply}, nil
}

func (s *helloServiceServer) LotsOfReplies(request *pb.HelloRequest, stream pb.HelloService_LotsOfRepliesServer) error {
	response := &pb.HelloResponse{Reply: s.reply}
	for i := 0; i < s.count; i++ {
		if err := stream.Send(response); err != nil {
			return err
		}
	}
	return nil
}

func (s *helloServiceServer) LotsOfGreetings(stream pb.HelloService_LotsOfGreetingsServer) error {
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloResponse{Reply: s.reply})
		}
		if err != nil {
			return err
		}
	}
}

func (s *helloServiceServer) BidiHello(stream pb.HelloService_BidiHelloServer) error {
	response := &pb.HelloResponse{Reply: s.reply}
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}

func newServer() *helloServiceServer {
	s := new(helloServiceServer)
	s.count = *count
	s.reply = *reply
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterHelloServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
