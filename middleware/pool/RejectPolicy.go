package pool

import (
	"fmt"
	"github.com/petermattis/goid"
)

type RejectPolicy interface {
	reject(task *TaskFuture) *TaskFuture
}

type rejectPolicyImpl func(task *TaskFuture) *TaskFuture

func (f rejectPolicyImpl) reject(task *TaskFuture) *TaskFuture {
	return f(task)
}

func CallerRunPolicy() rejectPolicyImpl {
	return func(task *TaskFuture) *TaskFuture {
		fmt.Println("CallerRunPolicy go id ", goid.Get())
		task.run()
		fmt.Println("RunFinish ", goid.Get())
		return task
	}
}

func DiscardPolicy() rejectPolicyImpl {
	return func(task *TaskFuture) *TaskFuture {
		fmt.Println("Discard Policy")
		return nil
	}
}
