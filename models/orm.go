package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"go-gin-mvc/utils"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var (
	dbm         *xorm.Engine
	dbs         []*xorm.Engine
	check_count int
)

func init() {
	var err error

	//主库添加
	m_dsn := utils.Config.Section("mysql_master").Key("master").String()

	dbm, err = xorm.NewEngine("mysql", m_dsn)
	dbm.SetMaxIdleConns(10)
	dbm.SetMaxOpenConns(200)
	dbm.ShowSQL(true)
	dbm.ShowExecTime(true)



	if err != nil {
		fmt.Printf("Fail to connect to master: %v", err)
		os.Exit(1)
	}

	//从库添加
	slaves := utils.Config.Section("mysql_slave").Keys()
	for _, s_dsn := range slaves {
		_dbs, err := xorm.NewEngine("mysql", s_dsn.String())
		_dbs.SetMaxIdleConns(10)
		_dbs.SetMaxOpenConns(200)
		_dbs.ShowSQL(true)
		_dbs.ShowExecTime(true)

		if err!=nil {
			fmt.Println(err)
		}else{
			dbs = append(dbs, _dbs)
		}
	}
}

func GetMaster() *xorm.Engine {
	return dbm
}

func GetSlave() *xorm.Engine {
	rand.Seed(time.Now().Unix())
	rn := rand.Intn(len(dbs) - 1)
	return dbs[rn]
}


func DbCheck() {
	check_count++
	fmt.Printf("Begin->数据库检查:第%d次\n",check_count)

	if dbm != nil {

		// Raw SQL
		dbm_err:= dbm.Ping()

		if dbm_err != nil {
			fmt.Println("=!!!=主库报警处理")
			fmt.Println(dbm_err)

		} else {
			fmt.Println("--主数据库查询正常！")
		}
	} else {
		fmt.Println("=!!!=主数据库连接异常！")
	}



	for i := 0; i < len(dbs); i++ {
		if dbs[i] != nil {

			// Raw SQL
			dbs_err := dbs[i].Ping()
			if dbs_err != nil {

				fmt.Printf("=!!!=从数据库 %d 查询异常\n", i)
				fmt.Println(dbs_err)

			} else {
				fmt.Printf("--从数据库 %d 查询正常\n", i)
			}
		} else {
			fmt.Println(strconv.Itoa(len(dbs)) + "==从数据库没连接")
		}
	}


	fmt.Println("==\n")


}