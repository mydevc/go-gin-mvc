package main

import (
	"go-gin-mvc/jobs"
	"go-gin-mvc/queue"
	_ "net/http/pprof"
)


func main() {
	subscriber := new(jobs.Subscribe)
	forever := make(chan bool)


	q := queue.NewQueue()

	//队列执行的任务需要注册方可执行
	q.PushJob("Dosome",jobs.HandlerFunc(subscriber.Dosome))
	//q.PushJob("Fusome",jobs.HandlerFunc(subscriber.Fusome))

	//提前规划好队列，可按延时时间来划分。可多个任务由一个队列来执行，也可以一个任务一个队列，一个队列可启动多个消费者
	go q.NewShareQueue("SomeQueue")
	//go q.NewShareQueue("SomeQueue")
	//go q.NewShareQueue("SomeQueue")


	defer q.Close()
	<-forever
}