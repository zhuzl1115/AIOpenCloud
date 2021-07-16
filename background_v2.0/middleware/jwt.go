package middleware

import (
	"background_v1.0/models"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey 		 string = "zhu ziling"
)

// JWT authorization

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var JWTres models.Result
		JWTres.Status = "0"
		token := c.Request.Header.Get("token")
		if token == "" {
			JWTres.Message = "no token"
			c.JSON(http.StatusOK, JWTres)
			c.Abort()
			return
		}

		log.Print("get token: ", token)
		j := NewJWT()
		// token parser
		claims, err := j.ParserToken(token)

		fmt.Println(claims, err)
		if err != nil {
			// check if token is expired
			if err == TokenExpired {
				JWTres.Message = "token is expired, please login again"
				c.JSON(http.StatusOK, JWTres)
				c.Abort()
				return
			}
			// other error
			JWTres.Message = err.Error()
			c.JSON(http.StatusOK, JWTres)
			c.Abort()
			return
		}
		// set claim as the middleware value to controller
		c.Set("claims", claims)

	}
}

// JWT struct with SignKey
type JWT struct {
	SigningKey []byte
}

// define claims

type CustomClaims struct {
	jwt.StandardClaims
	Name  string `json:"userName"`
	Email string `json:"email"`
}

// initialize JWT struct

func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// get SignKey

func GetSignKey() string {
	return SignKey
}

// create token SHA256 SignKey

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	// https://gowalker.org/github.com/dgrijalva/jwt-go#Token
	// return *Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// token parser and return the situation could not handle

func (j *JWT) ParserToken(tokenString string) (*CustomClaims, error) {
	// https://gowalker.org/github.com/dgrijalva/jwt-go#ParseWithClaims
	// func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	fmt.Println(token, err)
	if err != nil {
		// https://gowalker.org/github.com/dgrijalva/jwt-go#ValidationError
		// jwt.ValidationError
		if ve, ok := err.(*jwt.ValidationError); ok {
			// ValidationErrorMalformed: incorrect format
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
				// ValidationErrorExpired: expired token
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
				// ValidationErrorNotValidYet:invalid token
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}

		}
	}

	// check if token.claims is valid and token.claims2*CustomClaims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid

}

