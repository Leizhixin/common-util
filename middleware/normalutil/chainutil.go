package normalutil

type ChainResultWrapper[R interface{}] struct {
	Result  R
	EndFlag bool
}

func ExecuteChain[T interface{}, R interface{}](executors []T, matcher func(T, R) bool, execFunc func(T, R) *ChainResultWrapper[R]) R {
	var result R
	for _, executor := range executors {
		if matcher(executor, result) {
			wrapper := execFunc(executor, result)
			if wrapper == nil {
				break
			}
			result = wrapper.Result
			if wrapper.EndFlag {
				break
			}
		}
	}
	return result
}
