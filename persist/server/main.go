package main

import (
	"crawler-distributed/config"
	"crawler-distributed/persist/service"
	"crawler-distributed/support/grpcsupport"
	"flag"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"log"
)

func main() {
	// 使用从命令行传入的参数启动 Grpc 服务端
	port := flag.Int("port", 0, "the grpcsupport server port")
	flag.Parse()
	if *port == 0 {
		*port = 8800
	}

	StartItemSaverServer(fmt.Sprintf("127.0.0.1:%d", *port), config.ElasticSearchIndexWithGrpc)
}

// StartItemSaverServer 启动 ItemSaver Grpc 服务端
func StartItemSaverServer(address, index string) {
	log.Println("💫ItemSaver server is running...")

	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.ElasticSearchAddr},
	})
	if err != nil {
		log.Fatalf("[startServerGrpc] error create elasticsearch client: %v", err)
	}

	grpcsupport.NewGrpcItemSaverServer(config.Network, address, service.NewItemServer(esClient, index))
}
