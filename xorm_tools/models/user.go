package models

import "time"

type User struct {
	Age       int    `json:"age" xorm:"not null INT(11)"`
	CreatedAt time.Time    `json:"created_at" xorm:"created not null INT(11)"`
	Name      string `json:"name" xorm:"not null VARCHAR(20)"`
	UpdatedAt time.Time    `json:"updated_at" xorm:"updated not null INT(11)"`
	UserId    int    `json:"user_id" xorm:"not null pk autoincr comment('用户id') INT(11)"`
}
