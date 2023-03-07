package scheduler

import (
	"crawler-distributed/model"
)

// Scheduler 调度器
type Scheduler interface {
	// Submit 将待处理的 Request 请求送给调度器
	Submit(request model.Request)

	// CreateWorker 创建 Worker
	CreateWorker() chan model.Request

	// WorkerReady 表示已经有 worker 处于就绪状态，可以继续接收任务了
	WorkerReady(chan model.Request)

	Run()
}

// WorkerReadyNotify 告知调度器已经有处于空闲状态的 Worker 就绪了
// 在 createWorker 函数中需要使用到 WorkerReady 函数的功能，但在参数中将 Scheduler 全部传入过于繁重，因此将该功能单独提取出来
type WorkerReadyNotify interface {
	WorkerReady(chan model.Request)
}
