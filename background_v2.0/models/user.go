package models

import (
	"background_v1.0/dao"
	"fmt"
	"github.com/jinzhu/gorm"
)

// define user struct

type User struct {
	gorm.Model
	Name        string `json:"name"`
	Pwd         string `json:"password"`
	Mobile      string  `json:"mobile" gorm:"DEFAULT:0"`
	Email       string `json:"email" gorm:"type:varchar(40);unique_index;"` //set unique value
}

// define login requirement

type LoginReq struct {
	Name string `json:"name"`
	Pwd  string `json:"password"`
}

// user register

func Register(username, pwd string, mobile string, email string) error {
	fmt.Println(username, pwd, mobile, email)

	if CheckUser(username) {
		return fmt.Errorf("user already existed")
	}

	user := User{
		Name:  username,
		Pwd:   pwd,
		Mobile: mobile,
		Email: email,
	}
	insertErr := dao.DB.Model(&User{}).Create(&user).Error
	return insertErr
}

// check user information

func CheckUser(username string) bool {

	result := false
	var user User

	dbResult := dao.DB.Where("name = ?", username).Find(&user)
	if dbResult.Error != nil {
		fmt.Printf("user information error:%v\n", dbResult.Error)
	} else {
		result = true
	}
	return result
}

// check login information

func LoginCheck(login LoginReq) (bool, User, error) {
	userData := User{}
	userExist := false

	var user User
	dbErr := dao.DB.Where("name = ?", login.Name).Find(&user).Error

	if dbErr != nil {
		return userExist, userData, dbErr
	}
	if login.Name == user.Name && login.Pwd == user.Pwd {
		userExist = true
		userData.Name = user.Name
		userData.Email = user.Email
	}

	if !userExist {
		return userExist, userData, fmt.Errorf("password mistake")
	}
	return userExist, userData, nil
}


