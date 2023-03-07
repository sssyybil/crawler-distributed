package main

import (
	"crawler-distributed/config"
	"crawler-distributed/support/grpcsupport"
	"crawler-distributed/worker/service"
	"flag"
	"fmt"
	"log"
)

var port = flag.Int("port", 0, "the port for me to listen on")

func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Println("💔 must specify a worker server port")
	}
	log.Printf("🌛worker server is running...")

	grpcsupport.NewGrpcWorkerServer(
		config.Network,
		fmt.Sprintf("127.0.0.1:%d", *port),
		service.NewCrawlService(),
	)
}
