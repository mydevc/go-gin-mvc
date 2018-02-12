package queue_jobs

import "fmt"

var jobs map[string]JobFunc

type JobFunc func(json string) bool


//在这里注册
func init() {
	jobs = make(map[string]JobFunc)
	jobs["send_mail"] = send_mail
}

func CallJobFun(func_name string, json string) bool {
	 _, ok := jobs[func_name]
	if ok {
		return jobs[func_name](json)
	}
	return false
}

//测试用例
func send_mail(mail string) bool {
	fmt.Println(mail)
	return true
}
