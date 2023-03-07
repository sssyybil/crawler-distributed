package grpcsupport

import (
	"crawler-distributed/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// TODO Grpc 中是否无法创建通用的 client 呢？

func NewItemSaverClient(address string) pb.ItemSaverServiceClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("[support.NewItemSaverClient] cannot dial grpcsupport server: ", err)
	}

	return pb.NewItemSaverServiceClient(conn)
}

// CreateClientPool 根据给定的服务器地址，创建客户端连接
func CreateClientPool(hosts []string) chan *pb.ItemSaverServiceClient {
	var clients []*pb.ItemSaverServiceClient
	for _, h := range hosts {
		grpcClient := NewItemSaverClient(h)
		clients = append(clients, &grpcClient)

		log.Printf("ItemSaver Client connect to %s", h)
	}

	out := make(chan *pb.ItemSaverServiceClient)
	go func() {
		for {
			// 轮流分发 client，且每轮分发结束后继续下一轮分发
			for _, c := range clients {
				out <- c
			}
		}
	}()
	return out
}

func NewWorkerClient(address string) pb.CrawlServiceClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("[support.NewWorkerClient] cannot dial grpcsupport server: ", err)
	}
	return pb.NewCrawlServiceClient(conn)
}

// CreateWorkerClientPool 根据服务器地址创建多个 worker 客户端连接
func CreateWorkerClientPool(hosts []string) *chan pb.CrawlServiceClient {
	var clients []pb.CrawlServiceClient
	for _, h := range hosts {
		grpcClient := NewWorkerClient(h)
		clients = append(clients, grpcClient)
	}

	out := make(chan pb.CrawlServiceClient)
	go func() {
		clientCount := 1
		for {
			// 轮流分发 client，且每轮分发结束后继续下一轮分发
			for _, c := range clients {
				out <- c
				log.Printf("当前分发第 %d 个 worker client.\n", clientCount)
				clientCount++
			}
		}
	}()
	return &out
}
