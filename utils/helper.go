package utils

import (
	"bytes"
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"time"
	"log"
)


var RedisPool *redis.Pool

func init() {

	max_idle, _ := Config.Section("redis").Key("max_idle").Int()
	max_active, _ := Config.Section("redis").Key("max_active").Int()
	host := Config.Section("redis").Key("host").String()
	database, _ := Config.Section("redis").Key("database").Int()

	idle_timeout_int, _ := Config.Section("redis").Key("idle_timeout").Int64()
	idle_timeout := time.Duration(idle_timeout_int) * time.Second

	timeout_int, _ := Config.Section("redis").Key("timeout").Int64()
	timeout := time.Duration(timeout_int) * time.Second

	// 建立连接池
	RedisPool = &redis.Pool{
		MaxIdle:     max_idle,
		MaxActive:   max_active,
		IdleTimeout: idle_timeout,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", host,
				//redis.DialPassword(conf["Password"].(string)),
				redis.DialDatabase(database),
				redis.DialConnectTimeout(timeout),
				redis.DialReadTimeout(timeout),
				redis.DialWriteTimeout(timeout))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}
}


func SetCookie(c *gin.Context,name string, value string, maxAge int)  {

	domain:= Config.Section("cookie").Key("domain").String()

	c.SetCookie(name,value,maxAge,"/",domain,false,true)
}


func ByteEncoder(s interface{}) []byte {
	var enc_result bytes.Buffer
	enc := gob.NewEncoder(&enc_result)
	if err := enc.Encode(s); err != nil {
		log.Fatal("encode error:", err)
	}

	return enc_result.Bytes()
}
