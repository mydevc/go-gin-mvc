package queue_jobs

import "fmt"

var jobs map[string]JobFunc

type JobFunc func(json string) error


//在这里注册
func init() {
	jobs = make(map[string]JobFunc)
	jobs["send_mail"] = send_mail
}

func CallJobFun(func_name string, json string) error {
	return jobs[func_name](json)

}

//测试用例
func send_mail(mail string) error {
	fmt.Println(mail)
	return nil
}
