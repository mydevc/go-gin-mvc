package common

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
)

var redisconn redis.Conn

func initCache() {
	if redisconn==nil {
		conn, err := redis.Dial("tcp", GetConfig("redis", "address").String())
		if err == nil {
			redisconn = conn
			redisconn.Do("SELECT", 1)
		}else {
			fmt.Println("redis connect fail!")
		}

	}
}
func PutCache(key string, value string, timeout int) {
	if GetConfig("system","usecache").String() != "true" {
		fmt.Println("disable use cache")
		return
	}
	initCache()
	_,err :=redisconn.Do("SET", key,value,"EX", timeout )

	if err!=nil{
		fmt.Println(err)
	}

}
func GetCache(key string) string {
	if GetConfig("system","usecache").String() != "true" {
		fmt.Println("disable use cache")
		return ""
	}
	initCache()
	value,err := redis.String(redisconn.Do("GET",key))
	if err != nil {
		return ""
	}
	return value
}
