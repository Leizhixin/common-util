package normalutil

type MatchFunc[Handler interface{}] func(handler Handler, inputParams []any, lastResults []any) bool

type ExecuteFunc[Handler interface{}] func(handler Handler, inputParams []any, lastResults []any) ([]any, bool)

type ChainExecuteRequest[Handler interface{}] struct {
	executors   []Handler
	matchFunc   MatchFunc[Handler]
	execFunc    ExecuteFunc[Handler]
	inputParams []any
	results     []any
}

func NewChainExecuteRequest[Handler interface{}](executors []Handler, matchFunc MatchFunc[Handler],
	execFunc ExecuteFunc[Handler], inputParams ...any) *ChainExecuteRequest[Handler] {
	return &ChainExecuteRequest[Handler]{
		executors:   executors,
		matchFunc:   matchFunc,
		execFunc:    execFunc,
		inputParams: inputParams,
	}
}

func (c *ChainExecuteRequest[Handler]) Do() {
	var results []any
	for _, executor := range c.executors {
		if c.matchFunc(executor, c.inputParams, results) {
			curResults, isEnd := c.execFunc(executor, c.inputParams, results)
			results = curResults
			if isEnd {
				break
			}
		}
	}
	if results == nil {
		results = make([]any, 0)
	}
	c.results = results
}

func GetResult[Handler interface{}, RetType interface{}](c *ChainExecuteRequest[Handler],
	index int,
	defVal RetType) (RetType, bool) {
	if len(c.results) == 0 || len(c.results) <= index {
		return defVal, false
	}
	resVal := c.results[index]
	var retVal RetType
	var ok bool
	retVal, ok = resVal.(RetType)
	if !ok {
		return defVal, false
	}
	return retVal, true
}
