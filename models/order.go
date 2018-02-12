package models

import (
	"fmt"
	"gin_api/database"
)

type Order struct {
	OrderId     int `gorm:"primary_key"`
	MemberId    int
	OrderSn     string
	PaySn       string
	PayType     string
	OrderStatus int
	TotalPrice  float32
}

func GetOrders() (orders []Order, err error) {
	dbs := database.GetDB("master")

	//	var orders []Order
	var count int

	err = dbs.Offset(0 * 10).Limit(10).Where("order_id > ?", 59882).Find(&orders).Count(&count).Error

	if err != nil {
		// 错误处理...

		fmt.Println(err)
	}
	return orders, err
}
