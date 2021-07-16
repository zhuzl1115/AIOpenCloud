package models

/**
 * @Description
 * @Author 朱子凌
 * @Date 2021/7/8 16:24
 **/

//Model
type (
	BgImg struct {
		ID        int		`json:"id"`
		MD5ID	  string	`json:"md5id"`
		Order     string	`json:"order"`
		Show   	  string	`json:"show"`
		Type	  string	`json:"type"`
		CreatedAt string	`json:"created_at"`
		Path 	  string	`json:"path"`
		UserIP	  string	`json:"-"`
		Width 	  string	`json:"width"`
		Height	  string	`json:"height"`
}
	Result struct {
		Status	  string	`json:"status" gorm:"default:'-1'"`
		Message	  string	`json:"message"`
		Data      []BgImg	`json:"data"`
	}
)
//CRUD

func add(){

}