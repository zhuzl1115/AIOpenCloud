package background

import (
	"background_v1.0/dao"
	"background_v1.0/models"
	"bytes"
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

func Paginate(c *gin.Context/*r *http.Request*//*page int, pageSize int*/) func(db *gorm.DB) *gorm.DB {

	return func (db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.Query("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(c.Query("page_size"))
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
		fmt.Println(filenameOnly+"document name only!!!\n")

		//add unix time to md5
		unixTime := time.Now().Unix()
		unixTimeStr := fmt.Sprintf("%v", unixTime)
		var buffer bytes.Buffer
		buffer.WriteString(filenameOnly)
		buffer.WriteString(unixTimeStr)
		filenameTime := buffer.String()
		fmt.Println(filenameTime+"document name and time!!!\n")
		w := md5.New()
		io.WriteString(w, filenameTime)
		md5Str := fmt.Sprintf("%x", w.Sum(nil))

		//create dir to save files
		dir := fmt.Sprintf("%d-%02d-%02d",time.Now().Year(), time.Now().Month(),
			time.Now().Day())

		//check if the dir is existed
		if _, err := os.Stat(dir);os.IsNotExist(err){
			os.Mkdir(dir, os.ModePerm)
			//os.Mkdir(dir, os.ModePerm)
		}

		//get the destination
		dst := fmt.Sprintf("./%s/%s",dir, md5Str)


		// upload file to destination
		c.SaveUploadedFile(file, dst)

		//get current time
		unixTime =time.Now().Unix()
		timeStr := fmt.Sprintf("%v", unixTime)

		//create BgImg type struct
		BgImg1 := models.BgImg{
			MD5ID: md5Str,
			Order: strconv.Itoa(len(files)-index),
			Show: "1",
			Type: "default",
			CreatedAt: timeStr,
			Path: dst,
			UserIP: c.ClientIP(),
			Width: strconv.Itoa(width),
			Height: strconv.Itoa(height),
		}
		fmt.Println(dao.DB.NewRecord(&BgImg1))
		dao.DB.Debug().Create(&BgImg1)

		//add new created record to res
		res.Data = append(res.Data, BgImg1)
	}
	res.Message = fmt.Sprintf("'%d' uploaded!\n", len(files))
	c.JSON(http.StatusOK, res)

	//delete background
	var delres models.Result
	delres.Status = "0"
	DelFile := c.Query("task_id")
	var DelList models.BgImg

	//query the target to delete
	dao.DB.Debug().Where("md5_id = ?",DelFile).First(&DelList)

	//check if is the same user
	if DelList.UserIP == c.ClientIP() {
		dao.DB.Debug().Delete(&DelList)
		delres.Message = fmt.Sprintf("%s has been deleted!!", DelFile)
		delres.Data = append(delres.Data, DelList)
	} else {
		delres.Message = "delete failed!!"
	}
	delres.Status ="1"
	c.JSON(http.StatusOK, delres)

}
func UpdateOrder(c *gin.Context) {
	var orderres models.Result
	orderres.Status = "0"
	var UpdateList models.BgImg

	//query task_id and order from the url
	TaskID := c.Query("task_id")
	OrderUpdate := c.Query("order")

	//update the order in mysql
	err := dao.DB.Debug().Model(&UpdateList).Where("md5_id = ?",TaskID).Update("order",OrderUpdate).Error

	//check if the task_id is existed
	if err == nil{
		orderres.Message = "update finished!!"
	}else{
		orderres.Message = "could not find the task_id!!"
	}

	//return the ordered list by json format
	var OrderList []models.BgImg
	dao.DB.Debug().Order("`order`+0 asc").Find(&OrderList)
	orderres.Data = OrderList
	orderres.Status = "1"
	c.JSON(http.StatusOK, orderres)

}
func UpdateStatus(c *gin.Context) {
	var statusres models.Result
	statusres.Status = "0"
	var ShowList models.BgImg
	TaskID := c.Query("task_id")
	ShowUpdate := c.Query("show")
	err := dao.DB.Debug().Model(&ShowList).Where("md5_id = ?",TaskID).Update("Show", ShowUpdate).Error

	//check if the task_id is existed
	if err == nil{
		statusres.Message = "update finished!!"
	}else{
		statusres.Message = "could not find the task_id!!"
	}

	var List []models.BgImg
	dao.DB.Debug().Order("'order'+0 asc").Find(&List)
	statusres.Data = List
	statusres.Status = "1"
	c.JSON(http.StatusOK,statusres)
}
func GetList(c *gin.Context) {

	//get the list ordered by asc
	var OrderList []models.BgImg
	dao.DB.Debug().Order("`order`+0 asc").Scopes(Paginate(c)).Find(&OrderList)
	var res models.Result
	res.Status = "1"
	res.Data = OrderList
	c.JSON(http.StatusOK,&res)
}
