package jobs

type HandlerFunc func(args interface{}) error

func (f HandlerFunc) CallFunc(args interface{}) error { //是不是很熟悉的代码？
	return f(args)
}
