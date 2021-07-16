package usertoken

import (
	"background_v1.0/middleware"
	"background_v1.0/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

//define the information of register

type RegisterInfo struct {
	Mobile string  `json:"mobile"`
	Name  string `json:"name"`
	Pwd   string `json:"password"`
	Email string `json:"email"`
}

type  RegisterResult struct {
	Status string `json:"status"`
	Message string `json:"message"`
}
// interface for register

func RegisterUser(c *gin.Context) {
	var RegisInfo RegisterInfo
	var RegisRes RegisterResult
	bindErr := c.BindJSON(&RegisInfo)
	if bindErr == nil {
		// check if register is success and return the error
		err := models.Register(RegisInfo.Name, RegisInfo.Pwd, RegisInfo.Mobile, RegisInfo.Email)

		if err == nil {
			RegisRes.Status = "0"
			RegisRes.Message = "register success"
		} else {
			RegisRes.Status = "-1"
			RegisRes.Message = "register failed: "+err.Error()
		}
	} else {
		RegisRes.Status = "-1"
		RegisRes.Message = "Incorrect format of user information"
	}
	c.JSON(http.StatusOK, RegisRes)
}

// login result

type LoginData struct {
	Token string `json:"token"`
	Name string `json:"name"`
}

type LoginResult struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Data 	LoginData `json:"data"`
}

//login interface

func Login(c *gin.Context) {
	var LoginReq models.LoginReq
	var LoginRes LoginResult
	if c.ShouldBindJSON(&LoginReq) == nil {
		// check the login requirement
		isPass, user, err := models.LoginCheck(LoginReq)
		if isPass {
			generateToken(c, user, LoginRes)
		} else {
			LoginRes.Status = "-1"
			LoginRes.Message = err.Error()
			c.JSON(http.StatusOK, LoginRes)
		}

	} else {
		LoginRes.Status = "-1"
		LoginRes.Message = "incorrect user information format"
		c.JSON(http.StatusOK, LoginRes)
	}
}

//token generator
func generateToken(c *gin.Context, user models.User, LoginRes LoginResult) {
	// construct JWT instance
	j := middleware.NewJWT()

	// construct user claims
	claims := middleware.CustomClaims{
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 100), // effect time
			ExpiresAt: int64(time.Now().Unix() + 6000), // expiration time
			Issuer:    "zhu ziling",                    // jwt issuer
		},
	}

	// generate token
	token, err := j.CreateToken(claims)

	if err != nil {
		LoginRes.Status = "-1"
		LoginRes.Message = err.Error()
		c.JSON(http.StatusOK, LoginRes)
	}

	log.Println(token)
	// get user data
	data := LoginData{
		Name:  user.Name,
		Token: token,
	}

	//return login result by JSON
	LoginRes.Status = "0"
	LoginRes.Message = "login success"
	LoginRes.Data = data
	c.JSON(http.StatusOK, LoginRes)
	return
}
