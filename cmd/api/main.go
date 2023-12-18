package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ashiqsabith123/love-bytes-proto/match/pb"
	"github.com/ashiqsabith123/user-details-svc/pkg/config"
	"github.com/ashiqsabith123/user-details-svc/pkg/di"
	"github.com/ashiqsabith123/user-details-svc/pkg/helper"
	"google.golang.org/grpc"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(helper.Red("Error while loading config", err))
	}
	service := di.IntializeService(config)

	lis, err := net.Listen("tcp", config.Port.SvcPort)
	if err != nil {
		log.Fatalln("Failed to listening:", err)
	}

	fmt.Println("Match Svc on", config.Port.SvcPort)

	grpcServer := grpc.NewServer()

	pb.RegisterMatchServiceServer(grpcServer, &service)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc serve err: %v", err)
	}

}
