package main

import (
	"fmt"
	"net"

	logs "github.com/ashiqsabith123/love-bytes-proto/log"
	"github.com/ashiqsabith123/love-bytes-proto/match/pb"
	"github.com/ashiqsabith123/match-svc/pkg/config"
	"github.com/ashiqsabith123/match-svc/pkg/di"
	"google.golang.org/grpc"
)

func main() {

	// config, err := config.LoadConfig()
	// if err != nil {
	// 	log.Fatal(helper.Red("Error while loading config", err))
	// }
	// service := di.IntializeService(config)

	// lis, err := net.Listen("tcp", config.Port.SvcPort)
	// if err != nil {
	// 	log.Fatalln("Failed to listening:", err)
	// }

	// fmt.Println("Match Svc on", config.Port.SvcPort)

	// grpcServer := grpc.NewServer()

	// pb.RegisterMatchServiceServer(grpcServer, &service)

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Printf("Recovered from panic: %v", r)
	// 		// Additional handling or logging can be added here
	// 	}
	// }()

	// if err := grpcServer.Serve(lis); err != nil {
	// 	log.Fatalf("grpc serve err: %v", err)
	// }

	// Set up panic recovery

	StartService()

	// Start your microservice logic here

	// Sleep for a short duration before restarting

}

func StartService() {

	defer Recover()

	config, err := config.LoadConfig()

	if err != nil {
		logs.ErrLog.Fatalln("Error while loading config", err)
	}

	err = logs.InitLogger("./pkg/logs/log.log")
	if err != nil {
		fmt.Println(err)
		logs.ErrLog.Fatalln("Error while initilizing logger", err)
	}
	service := di.IntializeService(config)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v", r)
			// Additional handling or logging can be added here
		}
	}()

	lis, err := net.Listen("tcp", config.Port.SvcPort)
	if err != nil {
		logs.ErrLog.Fatalln("Failed to listening:", err)
	}

	logs.GenLog.Println("Match Svc listening on", config.Port.SvcPort)

	grpcServer := grpc.NewServer()

	pb.RegisterMatchServiceServer(grpcServer, &service)

	if err := grpcServer.Serve(lis); err != nil {
		logs.ErrLog.Fatalln("grpc serve err:", err)
	}

}

func Recover() {
	r := recover()

	fmt.Println("recovered from panic", r)
	StartService()
}
