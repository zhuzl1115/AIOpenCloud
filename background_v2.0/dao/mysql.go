package dao

import "github.com/jinzhu/gorm"

/**
 * @Description
 * @Author 朱子凌
 * @Date 2021/7/8 16:16
 **/

var (
	DB *gorm.DB
)

func InitMysql()(err error){
	DB, err = gorm.Open("mysql", "root:123456@(127.0.0.1:3306)/background?charset=utf8mb4&parseTime=True&loc=Local")
	if err !=nil{
		panic(err)
	}
	return DB.DB().Ping()
}

func Close()  {
	DB.Close()
}


