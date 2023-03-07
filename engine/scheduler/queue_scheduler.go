package scheduler

import (
	"crawler-distributed/model"
)

// QueuedScheduler 使用【队列】实现调度器
type QueuedScheduler struct {
	requestChan chan model.Request
	workerChan  chan chan model.Request // 存放 worker 的 channel，即每个 worker 都有自己的 channel
}

func (q *QueuedScheduler) Submit(request model.Request) {
	q.requestChan <- request
}

func (q *QueuedScheduler) CreateWorker() chan model.Request {
	return make(chan model.Request)
}

func (q *QueuedScheduler) WorkerReady(worker chan model.Request) {
	q.workerChan <- worker
}

// Run 将 request 和 worker 分别放入两个队列中，当队列中同时有 request 和 worker 在排队的情况下，再将 request 分发给 worker
func (q *QueuedScheduler) Run() {
	q.requestChan = make(chan model.Request)
	q.workerChan = make(chan chan model.Request)

	// 单独创建协程去做这件事的原因：
	go func() {
		var requestQueue []model.Request
		var workerQueue []chan model.Request

		for {
			var activeRequest model.Request
			var activeWorker chan model.Request

			// 当既有 request 在排队，又有 worker 在排队时，便可以将 request 分发给 worker
			if len(requestQueue) > 0 && len(workerQueue) > 0 {
				activeRequest = requestQueue[0]
				activeWorker = workerQueue[0]
			}

			select {
			case r := <-q.requestChan: // 当 requestChan 中有待处理的请求时，将请求存放到 requestQueue 中
				requestQueue = append(requestQueue, r)
			case w := <-q.workerChan: // 当 workerChan 中有空闲的 worker 时，将 worker 添加到 workerQueue 中
				workerQueue = append(workerQueue, w)
			case activeWorker <- activeRequest:
				workerQueue = workerQueue[1:]
				requestQueue = requestQueue[1:]
			}
		}
	}()
}
