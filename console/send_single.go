package main

import (
	"go-gin-mvc/jobs"
	"go-gin-mvc/queue"
	"strconv"
)

func main() {
	forever := make(chan bool)
	go func() {
		for i := 0; i < 1000000; i++ {
			queue.NewSender("SomeQueue", "Dosome", jobs.Subscribe{Name: "We are doing..." + strconv.Itoa(i)}).Send()
		}
	}()

	go func() {
		for i := 1000000; i < 1000000*2; i++ {
			queue.NewSender("SomeQueue", "Dosome", jobs.Subscribe{Name: "We are doing..." + strconv.Itoa(i)}).Send()
		}
	}()

	go func() {
		for i := 1000000 * 2; i < 1000000*3; i++ {
			queue.NewSender("SomeQueue", "Dosome", jobs.Subscribe{Name: "We are doing..." + strconv.Itoa(i)}).Send()
		}
	}()

	defer queue.SendConn.Close()
	<-forever

}
