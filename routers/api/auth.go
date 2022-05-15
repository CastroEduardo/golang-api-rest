package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CastroEduardo/golang-api-rest/conf"
	"github.com/CastroEduardo/golang-api-rest/models/authinterfaces"
	"github.com/CastroEduardo/golang-api-rest/pkg/app"
	"github.com/CastroEduardo/golang-api-rest/pkg/e"
	"github.com/CastroEduardo/golang-api-rest/pkg/logs_category"
	"github.com/CastroEduardo/golang-api-rest/pkg/setting"
	"github.com/CastroEduardo/golang-api-rest/pkg/util"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbsession_user_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbusers_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/logs_service"
	"github.com/astaxie/beego/validation"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember" valid:"Required;"`
}

type sendModel struct {
	Token   string `json:"token"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

type sendClaim struct {
	Data    interface{}
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
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
		return
	}
	//check UserFirts
	checkUser := dbusers_service.CheckUserPasswordForUser(modelUser.Username, modelUser.Password)

	if checkUser.NickName == "" {
		//validate user and password Email
		checkUser = dbusers_service.CheckUserPasswordForEmail(modelUser.Username, modelUser.Password)
	}

	if checkUser.NickName != "" {
		result := saveSessionUser(remember, checkUser)
		appG.Response(http.StatusOK, e.SUCCESS, result)
		return
	}

	userFailured(appG, username)

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
		return
	}
	//check UserFirts
	checkUser := dbusers_service.CheckUserPasswordForUser(modelUser.Username, modelUser.Password)

	if checkUser.NickName == "" {
		//validate user and password Email
		checkUser = dbusers_service.CheckUserPasswordForEmail(modelUser.Username, modelUser.Password)
	}

	if checkUser.NickName != "" {

		result := saveSessionUser(remember, checkUser)
		appG.Response(http.StatusOK, e.SUCCESS, result)
		return

	}

	userFailured(appG, username)

}

func userFailured(appG app.Gin, username string) {
	logs_service.Add(logs_category.USERFAILLOGIN, "TRY login FAILED FROM USER "+username, ipRequest)

	returnModel := sendModel{}
	returnModel.Success = false
	returnModel.Token = ""
	returnModel.Msg = conf.UserFailed

	fmt.Println(returnModel)
	appG.Response(http.StatusOK, e.SUCCESS, returnModel)
}

func saveSessionUser(remember bool, user authinterfaces.User) sendModel {

	var returnModel sendModel

	if user.Status == 0 {
		returnModel.Success = false
		returnModel.Token = ""
		returnModel.Msg = conf.UserDisabled

		logs_service.Add(logs_category.USERFAILLOGIN, "USER TRY LOGIN AND USER DISABLED "+user.NickName, ipRequest)

		return returnModel
	}

	tokenGet, err := util.GenerateToken("keep_trying _thief.", "keep_trying _thief.", remember)
	if err != nil {
		returnModel.Success = false
		returnModel.Token = ""
		returnModel.Msg = "error token"
		return returnModel
	}

	var timeLoggout time.Time
	if remember {
		timeLoggout = time.Now().Local().Add(time.Hour * time.Duration(setting.AppSetting.TimePersistToken))
	} else {
		timeLoggout = time.Now().Local().Add(time.Minute * time.Duration(setting.AppSetting.TimeNotPersistToken))
	}

	//add session //generate session User
	newSession := authinterfaces.SessionUser{
		Token:       tokenGet,
		Active:      true,
		DateAdd:     time.Now(),
		IdCompany:   user.IdCompany,
		IdUser:      user.ID,
		Remember:    remember,
		TokenExpire: timeLoggout,
	}
	dbsession_user_service.Add(newSession)

	logs_service.Add(logs_category.USERSLOGINSUCCESS, "USER LOGIN SUCCESS.. "+user.NickName, ipRequest)

	//fmt.Println(idSession)

	claim := dbsession_user_service.GetClaimForToken(tokenGet)

	resp := make(map[string]interface{})
	resp["company"] = claim.Company
	resp["user"] = claim.User
	resp["privileges"] = claim.UserPrivileges

	fmt.Println(resp["user"])
	//json.Marshal(claim.Company)
	//fmt.Println(string(jsonCompany))

	returnModel.Success = true
	returnModel.Token = tokenGet
	returnModel.Msg = conf.UserSuccess

	return returnModel
}

// @Summary Post ClaimUser
// @Produce  json
// @Tags Auth
//  //    @ID Authentication
// @Security bearerAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth/claim-user [post]
func PostClaimUser(c *gin.Context) {

	appG := app.Gin{C: c}
	ipRequest = c.ClientIP()
	//valid := validation.Validation{}

	sendData := sendClaim{}

	// token := c.Request.Header.Get("Authorization")

	// fmt.Println("Authorization: ", token)
	// c.JSON(200, gin.H{"Authorization": token})

	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		sendData.Success = false
		sendData.Msg = "Token Failed..."
		appG.Response(http.StatusOK, e.SUCCESS, sendData)
		return
	}

	//fmt.Println("Here")
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		sendData.Success = false
		sendData.Msg = "Token Failed..."
		appG.Response(http.StatusOK, e.SUCCESS, sendData)
		return
	}

	dataClaim := dbsession_user_service.GetClaimForToken(token)
	dataModel := make(map[string]interface{})

	dataModel["user"] = dataClaim.User
	dataModel["company"] = dataClaim.Company
	dataModel["privilege"] = dataClaim.UserPrivileges
	dataModel["rol"] = dataClaim.RolUser

	//fmt.Println(dataClaim)

	sendData.Success = true
	sendData.Msg = "Claim Success .."
	sendData.Data = "s" //dataClaim

	appG.Response(http.StatusOK, e.SUCCESS, dataModel)
}

// @Summary Post ClaimUser
// @Produce  json
// @Tags Auth
//  //    @ID Authentication
// @Security bearerAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth/logout [post]
func Postlogout(c *gin.Context) {

	appG := app.Gin{C: c}
	ipRequest = c.ClientIP()

	sendData := sendClaim{}

	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		sendData.Success = false
		sendData.Msg = "Token Failed..."
		appG.Response(http.StatusOK, e.SUCCESS, sendData)
		return
	}
	//fmt.Println("Here")
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		sendData.Success = false
		sendData.Msg = "Token Failed..."
		appG.Response(http.StatusOK, e.SUCCESS, sendData)
		return
	}

	logout := false
	session := dbsession_user_service.FindToToken(token)
	if session.Active {
		session.Active = false
	}
	result := dbsession_user_service.UpdateOne(session)
	if result {
		logout = true
	}
	appG.Response(http.StatusOK, e.SUCCESS, logout)
}

// @Summary Post ClaimUser
// @Produce  json
// @Tags Auth
// // @ID Authentication
// @Security bearerAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth/checkstatustoken [post]
func PostCheckStatusSession(c *gin.Context) {

	appG := app.Gin{C: c}
	ipRequest = c.ClientIP()

	var valid = false
	auth := c.Request.Header.Get("Authorization")

	if auth == "" {
		appG.Response(http.StatusOK, e.SUCCESS, valid)
		return
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		appG.Response(http.StatusOK, e.SUCCESS, valid)
		return
	}

	valid = true
	_, err := util.ParseToken(token)
	if err != nil {
		switch err.(*jwt.ValidationError).Errors {
		case jwt.ValidationErrorExpired:
			valid = false
		default:
			valid = false
		}
	}
	activeToken := dbsession_user_service.FindToToken(token)
	if !activeToken.Active {
		valid = false
	}
	appG.Response(http.StatusOK, e.SUCCESS, valid)

}

// // @Summary Get Auth
// // @Produce  json
// // @Tags Auth
// // @ID Authentication
// // @Param username query string true "username"
// // @Param password query string true "password"
// // @Param remember query bool true "remember"
// // @Success 200 {object} app.Response
// // @Failure 500 {object} app.Response
// // @Router /auth/getimg
// func Getimg(c *gin.Context) http.FileSystem {

// 	//appG := app.Gin{C: c}
// 	id := com.StrTo(c.Param("id")).MustInt()

// 	fmt.Println(id)

// 	//userFailured(appG, username)
// 	return http.Dir(upload.GetImageFullPath())
// }
