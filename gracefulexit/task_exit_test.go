package gracefulexit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type LongTimeWorkerTask struct {
	ctx    context.Context
	cancel func()
	wg     sync.WaitGroup
}

func (t *LongTimeWorkerTask) Start() {
	t.wg.Add(1)
	defer t.wg.Done()
	for {
		t.processTaskOnce(t.ctx)
		select {
		case <-t.ctx.Done():
			fmt.Println("Receive ctx done exist...")
			return
		case <-time.After(10 * time.Second): // 一个常驻内存的任务可能需要wakeup 处理
			fmt.Println("Wake up and reprocess")
			continue
		}
	}
}

// 如果处理一条数据的时间比较长 如阻塞的方式从redis里面拉取数据， 需要根据ctx是否取消进行处理
func (t *LongTimeWorkerTask) processTaskOnce(ctx context.Context) {

	select {
	case <-ctx.Done():
		fmt.Println("Going to exist")
	case <-time.After(5 * time.Second): // 模拟长时间执行的任务
		fmt.Println("Task Processing")
	}
}

func (t *LongTimeWorkerTask) Stop() {
	t.cancel()
	t.wg.Wait()
}

func ExampleLongTimeWorkerTask() {
	t := new(LongTimeWorkerTask)
	t.ctx, t.cancel = context.WithCancel(context.Background())

	go t.Start()
	time.Sleep(2 * time.Second)
	t.Stop()
	//output: Going to exist
	//Receive ctx done exist...
}
