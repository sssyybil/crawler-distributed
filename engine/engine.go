package engine

import (
	"context"
	"crawler-distributed/engine/scheduler"
	"crawler-distributed/model"
	"crawler-distributed/support/redissupport"
	"crawler-distributed/worker/service"
	"log"
)

type Engine struct {
	Scheduler        scheduler.Scheduler
	WorkerCount      int // Worker 服务数量
	ItemChan         chan model.Item
	RequestProcessor model.Processor
}

func (e *Engine) CreateEngineWorker(in chan model.Request, out chan model.ParserResult, ready scheduler.WorkerReadyNotify) {
	go func() {
		for {
			// 将空闲的 worker（worker 的类型就是存放了 Request 的 channel）交给调度器
			ready.WorkerReady(in)
			request := <-in
			parserResult, err := e.RequestProcessor(request)

			if err != nil {
				continue
			}
			out <- parserResult
		}
	}()
}

func (e *Engine) Run(ctx context.Context, seeds ...model.Request) {
	// 启动调度器
	e.Scheduler.Run()

	// 创建 Worker
	out := make(chan model.ParserResult)
	for i := 0; i < e.WorkerCount; i++ {
		e.CreateEngineWorker(e.Scheduler.CreateWorker(), out, e.Scheduler)
	}

	redisClient := redissupport.NewRedisClient()

	// 遍历种子，根据种子中的 URL 剔除掉已经存在的种子，将未被剔除掉的种子送入调度器中
	for _, r := range seeds {
		// URL 去重
		//if service.IsDuplicate(r.Url) {
		//	log.Printf("⚡️Duplicate request: %s", r.Url)
		//	continue
		//}
		if service.ReduplicateWithRedis(ctx, redisClient, r.Url) {
			log.Printf("⚡️Duplicate request: %s", r.Url)
			continue
		}
		e.Scheduler.Submit(r)
	}

	for {
		result := <-out
		// 从 worker 中接收的数据将分为两部分进行处理：有效的信息存入 ElasticSearch 中，请求信息再次送入调度器中进行解析
		for _, item := range result.Items {
			log.Printf("Engine Run result item %v\n", item)

			itemCopy := item
			go func() {
				e.ItemChan <- itemCopy
			}()
		}
		for _, r := range result.Requests {
			//if service.IsDuplicate(r.Url) {
			//	log.Printf("⚡️Duplicate request: %s", r.Url)
			//	continue
			//}
			if service.ReduplicateWithRedis(ctx, redisClient, r.Url) {
				log.Printf("⚡️Duplicate request: %s", r.Url)
				continue
			}
			e.Scheduler.Submit(r)
		}
	}
}
