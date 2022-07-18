package main

import (
	"context"
	"fmt"
	"github.com/arthur1470/fc2-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Pedro",
		Email: "j@j.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Pedro",
		Email: "j@j.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}

		fmt.Println("Status: ", stream.Status)
	}

	fmt.Println(responseStream)
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "w1",
			Name:  "Wesley",
			Email: "wes@wes.com",
		},
		&pb.User{
			Id:    "a1",
			Name:  "Arthur",
			Email: "art@art.com",
		},
		&pb.User{
			Id:    "a2",
			Name:  "Allan",
			Email: "all@all.com",
		},
		&pb.User{
			Id:    "r1",
			Name:  "Ronaldo",
			Email: "ron@ron.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "w1",
			Name:  "Wesley",
			Email: "wes@wes.com",
		},
		&pb.User{
			Id:    "a1",
			Name:  "Arthur",
			Email: "art@art.com",
		},
		&pb.User{
			Id:    "a2",
			Name:  "Allan",
			Email: "all@all.com",
		},
		&pb.User{
			Id:    "r1",
			Name:  "Ronaldo",
			Email: "ron@ron.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			//time.Sleep(time.Second * 2)
		}

		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			fmt.Println("")

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}

			fmt.Printf("Recebendo user %v com status: %v \n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
