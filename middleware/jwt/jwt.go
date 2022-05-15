package jwt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/CastroEduardo/golang-api-rest/pkg/e"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbsession_user_service"

	"github.com/CastroEduardo/golang-api-rest/pkg/util"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		//c.Request.URL.Path = "http://sgoogle.com"
		// oldpath := c.Request.URL.Path
		// result := strings.Replace(oldpath, "folder1", "766612", -1)
		// c.Request.URL.Path = result
		// // fmt.Println(c.Request.Host)
		// fmt.Println(c.Request.URL.Path)

		//Proted files in folder images####
		token_key := c.Query("token") //com.StrTo(c.Param("id")
		if token_key != "" {
			//not permit =/upload/images/ only
			longitud := len(c.Request.URL.Path)
			if longitud < 22 {
				c.Request.URL.Path = ""
			}
			folderCompany := "766612"
			//check min lent of folder to company
			longitud2 := len(c.Request.URL.Path + folderCompany)
			if longitud2 < 38 {
				c.Request.URL.Path = ""
			}
			isLogin := dbsession_user_service.FindToToken(token_key)
			if isLogin.Active {
				c.Next()
				return
			}
		}

		//fmt.Println(c.Request.RequestURI)
		// token2 := strings.Split(c.Request.Header["Authorization"][0], "dasd ")[1]
		// fmt.Println(token2)
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.String(http.StatusForbidden, "No Authorization header provided")
			c.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			c.String(http.StatusForbidden, "Could not find bearer token in Authorization header")
			c.Abort()
			return
		}
		fmt.Println("###########-->")
		var code int
		var data interface{}

		code = e.SUCCESS
		//token := c.Query("token")

		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			_, err := util.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				default:
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				}
			}

		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		// if !checkTokenDbSession(token) {
		// 	//if session is active
		// 	session := dbsession_user_service.FindToToken(token)

		// 	if !session.Active {
		// 		logs_service.Add(logs_category.FAILUREDTOKEN, "TRY ACCES FAILURED ..TOKEN DISABLED ==> IDSESSION : "+session.ID, "")
		// 		c.String(http.StatusRequestTimeout, "Session Expire...[ Check your token ]")
		// 		c.Abort()
		// 	}
		// }

		// fmt.Println(c.Request.RequestURI)
		// if c.Request.RequestURI == "/api/v2/upload/images/testimg.png" {
		// 	fmt.Println("__Envio")
		// 	c.Next()
		// 	return
		// }

		c.Next()
		return

	}
}

func checkTokenDbSession(token string) bool {
	return false
}
