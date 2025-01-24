package pool

import (
	"errors"
	"sync"
	"sync/atomic"
)

// RoutinePool 协程池
type RoutinePool struct {
	coreWorkerSize    int
	coreWorkerChan    chan struct{}
	maxWorkSize       int
	backEndWorkerChan chan struct{}
	taskQueue         chan Runnable
	wg                sync.WaitGroup
	closedChan        chan struct{}
	closed            atomic.Bool
	mainLock          sync.RWMutex
	rejectPolicy      RejectPolicy
}

func NewRoutingPool(coreSize, maxSize, queueSize int) *RoutinePool {
	rt := &RoutinePool{
		coreWorkerSize:    coreSize,
		coreWorkerChan:    make(chan struct{}, coreSize),
		maxWorkSize:       maxSize,
		backEndWorkerChan: make(chan struct{}, maxSize-coreSize),
		taskQueue:         make(chan Runnable, queueSize),
		wg:                sync.WaitGroup{},
		closedChan:        make(chan struct{}),
		closed:            atomic.Bool{},
	}
	rt.closed.Store(false)
	//rt.backTask()
	return rt
}

func NewRoutingPoolWithPolicy(coreSize, maxSize, queueSize int, rejectPolicy RejectPolicy) *RoutinePool {
	rt := &RoutinePool{
		coreWorkerSize:    coreSize,
		coreWorkerChan:    make(chan struct{}, coreSize),
		maxWorkSize:       maxSize,
		backEndWorkerChan: make(chan struct{}, maxSize-coreSize),
		taskQueue:         make(chan Runnable, queueSize),
		wg:                sync.WaitGroup{},
		closedChan:        make(chan struct{}),
		closed:            atomic.Bool{},
		rejectPolicy:      rejectPolicy,
	}
	rt.closed.Store(false)
	//rt.backTask()
	return rt
}

func (rt *RoutinePool) Submit(callable Callable) *TaskFuture {
	//fmt.Println("Submit go id ", goid.Get())
	if rt.closed.Load() {
		return nil
	}
	select {
	case <-rt.closedChan:
		return nil
	default:
	}
	rt.mainLock.RLock()
	defer rt.mainLock.RUnlock()
	if rt.closed.Load() {
		return nil
	}

	taskFuture := NewTaskFuture(callable)
	select {
	case rt.coreWorkerChan <- struct{}{}:
		//启动核心协程，Worker结束后一定记得在Worker协程中释放掉
		worker := NewWorker(rt, taskFuture)
		worker.Start()
		return taskFuture
	default:
	}
	select {
	//核心协程满了，插入任务队列
	case rt.taskQueue <- taskFuture:
		return taskFuture
	default:
	}
	//任务队列满了
	select {
	//先尝试增加Worker至最大Worker数量
	case rt.backEndWorkerChan <- struct{}{}:
		// Worker结束后一定记得释放掉
		worker := NewWorker(rt, nil)
		worker.Start()
		select {
		case rt.taskQueue <- taskFuture:
			//增加Worker成功后，再次尝试进入队列
			return taskFuture
		default:
			//增加Worker失败，直接返回结果失败，后续开发拒绝策略
			return rt.reject(taskFuture)
		}
	default:
		//增加Worker失败，直接返回结果失败
		return rt.reject(taskFuture)
	}
}

func (rt *RoutinePool) reject(task *TaskFuture) *TaskFuture {
	// 拒绝策略
	if rt.rejectPolicy != nil {
		return rt.rejectPolicy.reject(task)
	}
	return nil
}

func (rt *RoutinePool) Close() {
	if rt.closed.CompareAndSwap(false, true) {
		close(rt.closedChan)
		rt.mainLock.Lock()
		defer rt.mainLock.Unlock()
		close(rt.taskQueue)
		close(rt.coreWorkerChan)
		close(rt.backEndWorkerChan)
		rt.wg.Wait()
	}
}

type Runnable interface {
	run()
}

type Callable func() (result interface{}, err error)

func (f Callable) call() (result interface{}, err error) {
	return f()
}

type TaskFuture struct {
	result     interface{}
	err        error
	finishChan chan struct{}
	callable   Callable
}

func NewTaskFuture(callable Callable) *TaskFuture {
	return &TaskFuture{
		result:     nil,
		err:        nil,
		finishChan: make(chan struct{}),
		callable:   callable,
	}
}

func (t *TaskFuture) run() {
	t.result, t.err = t.callable.call()
	close(t.finishChan)
}

func (t *TaskFuture) Get() (interface{}, error) {
	<-t.finishChan
	return t.result, t.err
}

func (t *TaskFuture) GetNow() (interface{}, error) {
	select {
	case <-t.finishChan:
		return t.result, t.err
	default:
		return nil, errors.New("unfinished")
	}
}
