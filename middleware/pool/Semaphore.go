package pool

type Semaphore struct {
	ch chan interface{}
}

func NewSemaphore(limit int) *Semaphore {
	if limit <= 0 {
		panic("limit must greater than or equal to 1")
	}
	return &Semaphore{ch: make(chan interface{}, limit)}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}

func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}
