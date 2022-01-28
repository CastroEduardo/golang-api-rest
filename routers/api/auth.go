package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/CastroEduardo/golang-api-rest/pkg/app"
	"github.com/CastroEduardo/golang-api-rest/pkg/e"
	"github.com/CastroEduardo/golang-api-rest/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
	Remember bool   `valid:"Required;"`
}

// @Summary Get Auth
// @Produce  json
// @Tags Auth
// @ID Authentication
// @Param username query string true "userName"
// @Param password query string true "password"
// @Param remember query bool true "remember"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth [post]
func GetAuth(c *gin.Context) {

	fmt.Println("---###### here!!")
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	username := c.PostForm("username")
	password := c.PostForm("password")
	strTo, _ := strconv.ParseBool(c.PostForm("remember"))
	remember := strTo

	a := auth{Username: username, Password: password, Remember: remember}
	ok, _ := valid.Valid(&a)
	fmt.Println(ok)

	// if !ok {
	// 	app.MarkErrors(valid.Errors)
	// 	appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	// 	return
	// }

	// authService := auth_service.Auth{Username: username, Password: password}
	// isExist, err := authService.Check()
	// if err != nil {
	// 	appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
	// 	return
	// }

	// if !isExist {
	// 	appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
	// 	return
	// }

	token, err := util.GenerateToken(username, password, remember)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	fmt.Println(http.StatusOK)

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}
