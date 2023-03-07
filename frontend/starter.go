package main

import (
	"crawler-distributed/frontend/controller"
	"log"
	"net/http"
)

func main() {
	// 为访问 css 文件，启动一个文件服务，也可访问 frontend/view 目录下的文件
	http.Handle("/", http.FileServer(http.Dir("frontend/view")))
	http.Handle("/search", controller.CreateSearchResultHandler("frontend/view/template.html"))

	err := http.ListenAndServe(":8888", nil)

	if err != nil {
		log.Fatalf("⚡️Web Server error: %v\n", err)
	}
}
