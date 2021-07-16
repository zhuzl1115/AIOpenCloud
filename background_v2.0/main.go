package main

import (
	"background_v1.0/controller/background"
	"background_v1.0/controller/usertoken"
	"background_v1.0/dao"
	"background_v1.0/middleware"
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

	dao.DB.AutoMigrate(&models.BgImg{}, &models.User{})
	//routers
	router := gin.Default()
	router.POST("/background", middleware.JWTAuth(), background.UploadAndDelete)
	router.PUT("/background/order", middleware.JWTAuth(), background.UpdateOrder)
	router.PUT("/background/status", middleware.JWTAuth(), background.UpdateStatus)
	router.GET("/backgrounds", middleware.JWTAuth(), background.GetList)
	router.POST("/register", usertoken.RegisterUser)
	router.POST("/login", usertoken.Login)

	router.Run(":8085")

}
