package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/CastroEduardo/golang-api-rest/models/authinterfaces"
	"github.com/CastroEduardo/golang-api-rest/pkg/app"
	"github.com/CastroEduardo/golang-api-rest/pkg/e"
	"github.com/CastroEduardo/golang-api-rest/pkg/logs_category"
	"github.com/CastroEduardo/golang-api-rest/pkg/util"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbsession_user_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbusers_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/logs_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember" valid:"Required;"`
}

type generateToken struct {
	Token   string
	Msg     string
	Success bool
}

var ipRequest = ""

// @Summary Post Auth
// @Produce  json
// @Tags Auth
// @ID Authentication
// @Param username query string true "username"
// @Param password query string true "password"
// @Param remember query bool true "remember"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth [post]
func PostAuth(c *gin.Context) {

	appG := app.Gin{C: c}
	valid := validation.Validation{}

	ipRequest = c.ClientIP()

	user := auth{}
	c.BindJSON(&user) //get params from Body

	username := user.Username //c.PostForm("username")
	password := user.Password //c.Query("password") //c.PostForm("password")
	//strTo, _ := strconv.ParseBool(c.PostForm("remember"))
	remember := user.Remember //strTo

	modelUser := auth{Username: username, Password: password, Remember: remember}
	ok, _ := valid.Valid(&modelUser)

	if !ok {
		appG.Response(http.StatusInternalServerError, e.INVALID_PARAMS, nil)
		// 	return
	}
	//check UserFirts
	checkUser := dbusers_service.CheckUserPasswordForUser(modelUser.Username, modelUser.Password)

	if checkUser.NickName == "" {
		//validate user and password Email
		checkUser = dbusers_service.CheckUserPasswordForEmail(modelUser.Username, modelUser.Password)
	}

	if checkUser.NickName != "" {

		result := saveSessionUser(remember, checkUser)

		if !result.Success {
			appG.Response(http.StatusInternalServerError, e.ERROR, result)
			return
		} else {
			appG.Response(http.StatusOK, e.SUCCESS, result)
			return
		}

		// 	return

	}

	logs_service.Add(logs_category.USERFAILLOGIN, "TRY login FAILED FROM USER "+username, ipRequest)

	appG.Response(http.StatusInternalServerError, e.ERROR, map[string]string{
		"Error:": "Invalid user...",
	})
}

// @Summary Get Auth
// @Produce  json
// @Tags Auth
// @ID Authentication
// @Param username query string true "username"
// @Param password query string true "password"
// @Param remember query bool true "remember"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth [get]
func GetAuth(c *gin.Context) {

	appG := app.Gin{C: c}
	valid := validation.Validation{}

	ipRequest = c.ClientIP()

	username := c.Query("username") //c.PostForm("username")
	password := c.Query("password") //c.PostForm("password")
	strTo, _ := strconv.ParseBool(c.Query("remember"))
	remember := strTo

	modelUser := auth{Username: username, Password: password, Remember: remember}
	ok, _ := valid.Valid(&modelUser)

	if !ok {
		appG.Response(http.StatusInternalServerError, e.INVALID_PARAMS, nil)
		// 	return
	}
	//check UserFirts
	checkUser := dbusers_service.CheckUserPasswordForUser(modelUser.Username, modelUser.Password)

	if checkUser.NickName == "" {
		//validate user and password Email
		checkUser = dbusers_service.CheckUserPasswordForEmail(modelUser.Username, modelUser.Password)
	}

	if checkUser.NickName != "" {

		result := saveSessionUser(remember, checkUser)

		if !result.Success {
			appG.Response(http.StatusInternalServerError, e.ERROR, result)
			return
		} else {
			appG.Response(http.StatusOK, e.SUCCESS, result)
			return
		}

	}

	logs_service.Add(logs_category.USERFAILLOGIN, "TRY login FAILED FROM USER "+username, ipRequest)

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

	// token, err := util.GenerateToken(username, password, remember)
	// if err != nil {
	// 	appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
	// 	return
	// }

	// fmt.Println(http.StatusOK)

	appG.Response(http.StatusInternalServerError, e.ERROR, map[string]string{
		"Error:": "Invalid user...",
	})
}

func saveSessionUser(remember bool, user authinterfaces.User) generateToken {

	var returnModel generateToken

	tokenGet, err := util.GenerateToken("keep_trying _thief.", "keep_trying _thief.", remember)
	if err != nil {

		returnModel.Success = false
		returnModel.Token = ""
		returnModel.Msg = "error token"

		return returnModel
	}

	if user.Status == 0 {
		returnModel.Success = false
		returnModel.Token = ""
		returnModel.Msg = "disabled User.. "

		logs_service.Add(logs_category.USERFAILLOGIN, "USER TRY LOGIN AND USER DISABLED "+user.NickName, ipRequest)

		return returnModel
	}

	returnModel.Success = true
	returnModel.Token = tokenGet
	returnModel.Msg = "Success Login User.. "

	//add session //generate session User
	newSession := authinterfaces.SessionUser{
		Token:     tokenGet,
		Active:    true,
		DateAdd:   time.Now(),
		IdCompany: user.IdCompany,
		IdUser:    user.ID,
		Remember:  remember,
	}

	getIdSession := dbsession_user_service.Add(newSession)

	logs_service.Add(logs_category.USERSLOGINSUCCESS, "USER LOGIN SUCCESS.. "+user.NickName, ipRequest)

	fmt.Println(getIdSession)

	//fmt.Println(tokenGet)
	//fmt.Println(idUser)

	return returnModel
}
