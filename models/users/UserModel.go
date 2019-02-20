package model_users

import (
	"fmt"
	"go-gin-mvc/entitys"
	"go-gin-mvc/models"
)



func UserAdd(user *entitys.User)  {

	orm:= models.GetMaster()

	_,err:=orm.Insert(user)

	if err != nil {
		fmt.Println(err)
		return
	}
}


func UserOne() (*entitys.User,bool) {

	orm:= models.GetSlave()

	user := new(entitys.User)



	has, err := orm.Where("name=?", "laodeng").Get(user)
	if err!=nil {
		fmt.Println(err)
	}
	return user,has
}

func UserList() ([]entitys.User,error) {

	orm:= models.GetSlave()

	users := make([]entitys.User, 0)

	err := orm.Where("name<>''").Find(&users)

	return users,err
}


