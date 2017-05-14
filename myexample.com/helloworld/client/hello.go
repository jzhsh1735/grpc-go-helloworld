package main

import (
	"flag"
	"io"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "myexample.com/helloworld/hello"
)

var (
	server = flag.String("server", "localhost:12345", "The server address")
	greeting = flag.String("greeting", "HelloServer!", "The greeting")
	count = flag.Int("count", 3, "The client greeting count")
)

func SayHello(client pb.HelloServiceClient) {
	grpclog.Printf("SayHello")
	response, err := client.SayHello(context.Background(), &pb.HelloRequest{Greeting: *greeting})
	if err != nil {
		grpclog.Fatalf("%v: %v", client, err)
	}
	grpclog.Println(response)
}

func LotsOfReplies(client pb.HelloServiceClient) {
	grpclog.Printf("LotsOfReplies")
	stream, err := client.LotsOfReplies(context.Background(), &pb.HelloRequest{Greeting: *greeting})
	if err != nil {
		grpclog.Fatalf("%v: %v", client, err)
	}
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			grpclog.Fatalf("%v: %v", client, err)
		}
		grpclog.Println(response)
	}
}

func LotsOfGreetings(client pb.HelloServiceClient) {
	grpclog.Printf("LotsOfGreetings")
	request := &pb.HelloRequest{Greeting: *greeting}
	stream, err := client.LotsOfGreetings(context.Background())
	if err != nil {
		grpclog.Fatalf("%v: %v", client, err)
	}
	for i := 0; i < *count; i++ {
		if err := stream.Send(request); err != nil {
			grpclog.Fatalf("%v: %v", client, err)
		}
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		grpclog.Fatalf("%v: %v", client, err)
	}
	grpclog.Println(response)
}

func BidiHello(client pb.HelloServiceClient) {
	grpclog.Printf("BidiHello")
	request := &pb.HelloRequest{Greeting: *greeting}
	stream, err := client.BidiHello(context.Background())
	if err != nil {
		grpclog.Fatalf("%v: %v", client, err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				grpclog.Fatalf("%v: %v", client, err)
			}
			grpclog.Println(response)
		}
	}()
	for i := 0; i < *count; i++ {
		if err := stream.Send(request); err != nil {
			grpclog.Fatalf("%v: %v", client, err)
		}
	}
	stream.CloseSend()
	<-waitc
}

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*server, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewHelloServiceClient(conn)
	SayHello(client)
	LotsOfReplies(client)
	LotsOfGreetings(client)
	BidiHello(client)
}
