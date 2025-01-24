package pool

import (
	"fmt"
	"github.com/petermattis/goid"
	"testing"
)

// 这回不会再有死锁了
func TestNewRoutingPool(t *testing.T) {
	rp := NewRoutingPoolWithPolicy(5, 10, 10, CallerRunPolicy())
	t.Cleanup(func() {
		rp.Close()
	})
	//rp := NewRoutingPool(5, 5, 10)
	fmt.Println("main go routine id ", goid.Get())
	//defer rp.Close()
	var taskList []*TaskFuture
	for i := 0; i <= 100; i++ {
		num := i
		taskFuture := rp.Submit(func() (result interface{}, err error) {
			fmt.Println("task go routine id ", goid.Get())
			result = num
			return
		})
		if taskFuture != nil {
			taskList = append(taskList, taskFuture)
		}
	}
	for _, e := range taskList {
		result, _ := e.Get()
		fmt.Println("task future result ", result)
	}
	fmt.Println("all finished")
}
