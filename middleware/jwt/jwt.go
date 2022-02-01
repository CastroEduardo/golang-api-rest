package jwt

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/CastroEduardo/golang-api-rest/pkg/e"
	"github.com/CastroEduardo/golang-api-rest/pkg/logs_category"

	"github.com/CastroEduardo/golang-api-rest/pkg/util"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbsession_user_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/logs_service"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {

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

		if !checkTokenDbSession(token) {
			//if session is active
			session := dbsession_user_service.FindToToken(token)
			if !session.Active {

				logs_service.Add(logs_category.FAILUREDTOKEN, "TRY ACCES FAILURED ..TOKEN DISABLED ==> IDSESSION : "+session.ID, "")

				c.String(http.StatusRequestTimeout, "Session Expire...[ Check your token ]")
				c.Abort()
			}

			return
		}

		c.Next()
	}
}

func checkTokenDbSession(token string) bool {

	return false
}
