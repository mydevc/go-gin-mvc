package queue

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/streadway/amqp"
	"go-gin-mvc/utils"
	"log"
)

type Sender struct {
	QueueName string //队列名称
	Job       string //任务名称
	Msg       []byte //消息体,回调时使用
}

var SendConn *amqp.Connection

func init(){
	if SendConn==nil {
		var err error
		SendConn, err = amqp.Dial("amqp://admin:dirdir@10.47.208.83:5672/")

		fmt.Println("New connect to RabbitMQ")

		failOnError(err, "Failed to connect to RabbitMQ")
	}
}



func NewSender(queue_name string, job string, send_msg interface{}) *Sender {
	msg:= utils.ByteEncoder(send_msg)

	//队列注册检查，可能需要借助redis来检查，暂放下
	return &Sender{QueueName: queue_name, Job: job, Msg: msg}
}


func (sender *Sender) Send() {
	var enc_result bytes.Buffer
	enc := gob.NewEncoder(&enc_result)
	if err := enc.Encode(sender); err != nil {
		log.Fatal("encode error:", err)
	}

	ch, err := SendConn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		sender.QueueName, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        enc_result.Bytes(),
		})
	//log.Printf(" [x] Sent %s", enc_result.Bytes())
	failOnError(err, "Failed to publish a message")

}
