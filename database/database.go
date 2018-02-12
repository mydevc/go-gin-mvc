package database

import (
	"fmt"
	"strconv"
	"github.com/jinzhu/gorm"
	"gin_api/common"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	dbKey       int
	DBM         *gorm.DB
	DBS         []*gorm.DB
	check_count int
)

func init() {

	GetDB("master")
	GetDB("slave")

	fmt.Println(DBS)
}

func GetDB(dbType string) *gorm.DB {

	if dbType == "master" {

		var err error

		if DBM == nil {
			dsn := common.GetConfig("mysql","masterDsn").String()
			DBM, err = DbConn(dsn)
			if err == nil {
				return DBM
			}
		}
		return DBM
	}

	if dbType == "slave" {

		//		fmt.Println("slave get conn")

		if len(DBS) < 1 {

			fmt.Println("slave new conn")

			slave_count, _ := common.GetConfig("mysql","slaveCount").Int()
			DBS = make([]*gorm.DB, slave_count)
			for i := 0; i < slave_count; i++ {

				dsn := "slaveDsn" + strconv.Itoa(i+1)
				dsn = common.GetConfig("mysql",dsn).String()

				fmt.Println("slave dns is : " + dsn)

				dbs, err := DbConn(dsn)
				if err == nil {
					DBS[i] = dbs
				} else {
					fmt.Println(err)
				}
			}
		}

		if len(DBS) > 0 {
			if dbKey > len(DBS)-1 || dbKey < 1 {
				dbKey = 0
			} else {
				dbKey++
			}

			return DBS[dbKey]
		}

	}

	return nil

}

func DbConn(dsn string) (*gorm.DB, error) {

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return db, err

	} else {
		fmt.Println("恭喜，数据库连接成功")
	}

		db.LogMode(true)
	// 全局禁用表名复数
	db.SingularTable(true) // 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "yf_" + defaultTableName
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(200)

	return db, err

}

func CheckDb() {
	fmt.Println("checkdb 执行...")

	for i := 0; i < 1; i++ {
		if DBM != nil {

			// Raw SQL
			rows, dbm_err := DBM.Raw("select 1 from mysql.db limit 1").Rows()

			if dbm_err != nil {
				fmt.Println("master fail ! 报警处理=================================")
				fmt.Println(dbm_err)

				//panic("db fail")

				//尝试重连接
				//GetDB("master")

			} else {
				defer rows.Close()
				fmt.Println(strconv.Itoa(check_count) + "--主数据库查询正常\n")
			}
		} else {
			fmt.Println(strconv.Itoa(check_count) + "主数据库没连接")
		}

		check_count++
	}

	for i := 0; i < len(DBS); i++ {
		if DBS[i] != nil {

			// Raw SQL
			rows, dbs_err := DBS[i].Raw("select 1 from mysql.db limit 1").Rows()
			if dbs_err != nil {
				fmt.Println("slave fail ! 报警处理")
				fmt.Println(dbs_err)

			} else {
				defer rows.Close()
				fmt.Printf("从数据库%d查询正常\n", i)
			}
		} else {
			fmt.Println(strconv.Itoa(len(DBS)) + "从数据库没连接")
		}
	}



}
