package background

import (
	"background_v1.0/dao"
	"background_v1.0/middleware"
	"background_v1.0/models"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

/**
 * @Description
 * @Author 朱子凌
 * @Date 2021/7/8 15:57
 **/

/**
 * @Description Pagination
 * @Param
 * @return
 **/

func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {

	return func (db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
func UploadAndDelete(c *gin.Context) {
	// upload multiple files
	var res models.Result
	res.Status = "0"
	form, _ :=c.MultipartForm()
	files, err := form.File["upload[]"]
	if err != true {
		res.Message = "uploading failed!"
		c.JSON(http.StatusOK, res)
		return
	}

	//create operation for multiple files
	for index, file := range files{
		log.Println(file.Filename)

		//get size of the image
		src, err := file.Open()
		if err != nil{
			fmt.Println("err = ", err)
			return
		}
		img ,_ ,err := image.DecodeConfig(src)
		if err != nil{
			fmt.Println("err = ", err)
			return
		}
		width := img.Width
		height := img.Height

		//filename2md5
		filenameSuffix := path.Ext(file.Filename)
		filenameOnly := strings.TrimSuffix(file.Filename, filenameSuffix)

		//add unix time to md5
		unixTimeStr := fmt.Sprintf("%v", time.Now().Unix())
		w := md5.New()
		io.WriteString(w, filenameOnly + unixTimeStr)
		md5Str := fmt.Sprintf("%x", w.Sum(nil))

		//create dir to save files
		dir := fmt.Sprintf("%d-%02d-%02d",time.Now().Year(), time.Now().Month(),
			time.Now().Day())

		//check if the dir is existed
		if _, err := os.Stat(dir);os.IsNotExist(err){
			os.Mkdir(dir, os.ModePerm)
		}

		//get the destination
		dst := fmt.Sprintf("./%s/%s",dir, md5Str)

		// upload file to destination
		c.SaveUploadedFile(file, dst)

		//get current time
		timeStr := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",time.Now().Year(), time.Now().Month(),
			time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())

		//create BgImg type struct
		BgImg1 := models.BgImg{
			MD5ID: md5Str,
			Order: strconv.Itoa(len(files)-index),
			Show: "1",
			Type: "default",
			CreatedAt: timeStr,
			Path: dst+filenameSuffix,
			UserIP: c.ClientIP(),
			Width: strconv.Itoa(width),
			Height: strconv.Itoa(height),
		}
		fmt.Println(dao.DB.NewRecord(&BgImg1))
		dao.DB.Debug().Create(&BgImg1)

		//add new created record to res
		res.Data = append(res.Data, BgImg1)
	}
	claims := c.MustGet("claims").(*middleware.CustomClaims)
	var ClaimsMessage string
	if claims != nil {
		ClaimsMessage = claims.Name+" token passed AND"
	}


	//delete background
	DelFile := c.Query("task_id")
	var DelList models.BgImg
	//query the target to delete
	dao.DB.Debug().Where("md5_id = ?",DelFile).First(&DelList)
	var DelMessage string
	//check if is the same user
	if DelList.UserIP == c.ClientIP() {
		dao.DB.Debug().Delete(&DelList)
		DelMessage = fmt.Sprintf("%s has been deleted!!", DelFile)
	} else {
		DelMessage = "delete failed!!"
	}
	res.Message = fmt.Sprintf("%s '%d' uploaded! %s", ClaimsMessage, len(files), DelMessage)
	c.JSON(http.StatusOK, res)
}

func UpdateOrder(c *gin.Context) {
	var OrderRes models.Result
	OrderRes.Status = "0"
	claims := c.MustGet("claims").(*middleware.CustomClaims)
	var ClaimsMessage string
	if claims != nil {
		ClaimsMessage = claims.Name+" token passed AND "
	}
	var UpdateList models.BgImg

	//query task_id and order from the url
	TaskID := c.Query("task_id")
	OrderUpdate := c.Query("order")

	//update the order in mysql
	err := dao.DB.Debug().Model(&UpdateList).Where("md5_id = ?",TaskID).Update("order",OrderUpdate).Error

	//check if the task_id is existed
	if err == nil{
		OrderRes.Message = ClaimsMessage + "update finished!!"
	}else{
		OrderRes.Message = ClaimsMessage + "could not find the task_id!!"
	}

	//return the ordered list by json format
	var OrderList []models.BgImg
	dao.DB.Debug().Order("`order`+0 asc").Find(&OrderList)
	OrderRes.Data = OrderList
	c.JSON(http.StatusOK, OrderRes)

}
func UpdateStatus(c *gin.Context) {
	var StatusRes models.Result
	StatusRes.Status = "0"
	claims := c.MustGet("claims").(*middleware.CustomClaims)
	var ClaimsMessage string
	if claims != nil {
		ClaimsMessage = claims.Name+" token passed AND "
	}
	var ShowList models.BgImg
	TaskID := c.Query("task_id")
	ShowUpdate := c.Query("show")
	err := dao.DB.Debug().Model(&ShowList).Where("md5_id = ?",TaskID).Update("Show", ShowUpdate).Error

	//check if the task_id is existed
	if err == nil{
		StatusRes.Message = ClaimsMessage + "update finished!!"
	}else{
		StatusRes.Message = ClaimsMessage + "could not find the task_id!!"
	}

	var List []models.BgImg
	dao.DB.Debug().Order("'order'+0 asc").Find(&List)
	StatusRes.Data = List
	c.JSON(http.StatusOK,StatusRes)
}
func GetList(c *gin.Context) {
	var res models.Result
	claims := c.MustGet("claims").(*middleware.CustomClaims)
	if claims != nil {
		res.Message = claims.Name+" token passed"
	}
	//get the list ordered by asc
	var OrderList []models.BgImg
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	dao.DB.Debug().Order("`order`+0 asc").Scopes(Paginate(page, pageSize)).Find(&OrderList)
	res.Status = "0"
	res.Data = OrderList
	c.JSON(http.StatusOK,&res)
}
