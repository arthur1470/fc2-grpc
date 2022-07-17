package main

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")

	if err != null {
		log.Fatalf("Could not connect: %v", err)
	}

	grpcServer := grpc.NewServer()

	if err := grpcServer.Serve(lis); err != null {
		log.Fatalf("Could not serve: %v", err)
	}

}