/*
  xcl (2015-8-15)
  多TubeName 多消费者
*/

package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/kr/beanstalk"
	"gin_api/console/queue_jobs"
)

var (
	TubeName1 string = "channel1"
	TubeName2 string = "channel2"
)

//添加任务 测试地址：http://localhost:8080/queue

func Consumer(fname, tubeName string) {
	if fname == "" || tubeName == "" {
		return
	}

	c, err := beanstalk.Dial("tcp", "192.168.1.168:11300")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c.Tube.Name = tubeName
	c.TubeSet.Name[tubeName] = true

	fmt.Println(fname, " [Consumer] tubeName:", tubeName, " c.Tube.Name:", c.Tube.Name)

	substr := "timeout"
	for {
		//fmt.Println(fname, " [Consumer]///////////////////////// begin..")

		//从队列中取出
		id, body, err := c.Reserve(3600*24*365*time.Second)

		fmt.Println("Reserve",string(body))

		if err != nil {
			if !strings.Contains(err.Error(), substr) {
				fmt.Println(fname, " [Consumer] [", c.Tube.Name, "] err:", err, " id:", id)
			}
			continue
		}
		//fmt.Println(fname, " [Consumer] [", c.Tube.Name, "] job:", id, " body:", string(body))

		call_result := queue_jobs.CallJobFun("send_mail",string(body))

		if call_result != nil {
			fmt.Println("调用出错了")
		}else {
			fmt.Println("调用完成")
		}

		//从队列中清掉
		err = c.Delete(id)


		//if err != nil {
		//	fmt.Println(fname, " [Consumer] [", c.Tube.Name, "] Delete err:", err, " id:", id)
		//} else {
		//	fmt.Println(fname, " [Consumer] [", c.Tube.Name, "] Successfully deleted. id:", id)
		//}
		//fmt.Println(fname, " [Consumer]///////////////////////// end..")
		//time.Sleep(1 * time.Second)
	}
	fmt.Println("Consumer() end. ")


}

var end chan string


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//Producer("PA", TubeName1)
	//Producer("PB", TubeName1)
	//Producer("PA", TubeName2)
	//Producer("PB", TubeName2)

	end = make(chan string)
	go Consumer("CA", TubeName1)
	go Consumer("CB", TubeName2)

	//time.Sleep(10 * time.Second)
	//for {
	//	time.Sleep(3600*time.Second)
	//}
	//close(end)

	for x := range end {
		fmt.Println(x)
		if x=="close"{
			//close(end)
		}
	}




}
