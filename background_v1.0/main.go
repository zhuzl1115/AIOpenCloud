package main

import (
	"background_v1.0/controller/background"
	"background_v1.0/dao"
	"background_v1.0/models"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
/**
 * @Description
 * @Author 朱子凌
 * @Date 2021/7/8 15:54
 **/

func main() {
	err :=dao.InitMysql()
	if err != nil{
		panic(err)
	}
	defer dao.Close()

	dao.DB.AutoMigrate(&models.BgImg{})
	//routers
	router := gin.Default()
	router.POST("/background", background.UploadAndDelete)
	router.PUT("/background/order", background.UpdateOrder)
	router.PUT("/background/status", background.UpdateStatus)
	router.GET("/backgrounds", background.GetList)


	router.Run(":8085")

}
