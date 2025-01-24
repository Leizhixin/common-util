package pool

type LimitTaskQueue chan *TaskFuture

func NewLimitTaskQueue(size int) LimitTaskQueue {
	return make(LimitTaskQueue, size)
}

func (q LimitTaskQueue) Append(future *TaskFuture) error {
	return q.AppendWithHandler(future, nil)
}

func (q LimitTaskQueue) AppendWithHandler(future *TaskFuture, resultHandler func(interface{})) error {
	select {
	case q <- future:
	default:
		err := q.FlushWithHandler(resultHandler)
		if err != nil {
			return err
		}
		q <- future
	}
	return nil
}

func (q LimitTaskQueue) Flush() error {
	return q.FlushWithHandler(nil)
}

func (q LimitTaskQueue) FlushWithHandler(resultHandler func(interface{})) error {
outer:
	for {
		select {
		case future, ok := <-q:
			if !ok {
				break outer
			}
			res, err := future.Get()
			if err != nil {
				return err
			}
			if resultHandler != nil {
				resultHandler(res)
			}
		default:
			break outer
		}
	}
	return nil
}
