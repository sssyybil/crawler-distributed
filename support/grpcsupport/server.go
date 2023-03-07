package grpcsupport

import (
	"crawler-distributed/pb"
	"crawler-distributed/persist/service"
	workerService "crawler-distributed/worker/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewGrpcItemSaverServer(network, address string, service *service.ItemSaverService) {
	grpcServer := grpc.NewServer()
	pb.RegisterItemSaverServiceServer(grpcServer, service)

	listener, err := net.Listen(network, address)
	if err != nil {
		log.Fatal("[support.NewGrpcItemSaverServer] cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("[support.NewGrpcItemSaverServer] cannot start grpcsupport server: ", err)
	}
}

func NewGrpcWorkerServer(network, address string, service *workerService.CrawlService) {
	grpcServer := grpc.NewServer()
	pb.RegisterCrawlServiceServer(grpcServer, service)

	listener, err := net.Listen(network, address)
	if err != nil {
		log.Fatal("[support.NewGrpcWorkerServer] cannot start server: ", err)
	}

	log.Printf("[support.NewGrpcWorkerServer] serer start to listening...\n")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("[support.NewGrpcWorkerServer] cannot start grpcsupport server: ", err)
	}
}
