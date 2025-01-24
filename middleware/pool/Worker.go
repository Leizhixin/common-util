package pool

type Worker struct {
	pool    *RoutinePool
	curTask *TaskFuture
}

func NewWorker(pool *RoutinePool, curTask *TaskFuture) *Worker {
	return &Worker{
		pool:    pool,
		curTask: curTask,
	}
}

func (w Worker) Start() {
	w.pool.wg.Add(1)
	go func() {
		isCoreWorker := false
		defer func() {
			// 当前worker销毁后需要释放对应的WorkerChannel 否则会有问题 这部分是必须执行的 所以需要被放到defer中
			if isCoreWorker {
				<-w.pool.coreWorkerChan
			} else {
				<-w.pool.backEndWorkerChan
			}
			w.pool.wg.Done()
		}()
		if w.curTask != nil {
			isCoreWorker = true
			w.curTask.run()
		}
	outer:
		for {
			select {
			case task, ok := <-w.pool.taskQueue:
				if !ok {
					break outer
				}
				task.run()
			default:
				break outer
			}
		}

	}()
}
