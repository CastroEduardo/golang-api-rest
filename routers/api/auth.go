package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CastroEduardo/golang-api-rest/conf"
	"github.com/CastroEduardo/golang-api-rest/models"
	"github.com/CastroEduardo/golang-api-rest/models/authinterfaces"
	"github.com/CastroEduardo/golang-api-rest/pkg/app"
	"github.com/CastroEduardo/golang-api-rest/pkg/e"
	"github.com/CastroEduardo/golang-api-rest/pkg/setting"
	"github.com/CastroEduardo/golang-api-rest/pkg/util"

	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbDevicesSys_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dblogs_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbsession_user_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbusers_service"
	"github.com/astaxie/beego/validation"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember" valid:"Required;"`
	From     string `json:"from"`
}

type deviceId struct {
	IdDevice string `json:"idDevice"`
	Type     string `json:"type"`
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
	fromLogin := user.From

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
		result := saveSessionUser(remember, checkUser, fromLogin)
		//fmt.Println(result)
		appG.Response(http.StatusOK, e.SUCCESS, result)
		return
	}

	userFailured(appG, username, fromLogin)

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
	fromLogin := c.Query("from")
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

		result := saveSessionUser(remember, checkUser, fromLogin)
		appG.Response(http.StatusOK, e.SUCCESS, result)
		return

	}

	userFailured(appG, username, fromLogin)

}

func userFailured(appG app.Gin, username string, fromLogin string) {

	dblogs_service.Add(conf.LOGIN_FAILE, "TRY login FAILED FROM USER "+username+" ( From : "+fromLogin+" )", "", ipRequest)

	returnModel := sendModel{}
	returnModel.Success = false
	returnModel.Token = ""
	returnModel.Msg = conf.UserFailed

	//fmt.Println(returnModel)
	appG.Response(http.StatusOK, e.SUCCESS, returnModel)
}

func saveSessionUser(remember bool, user authinterfaces.User_sys, fromLogin string) sendModel {

	var returnModel sendModel

	if user.Status == 0 {
		returnModel.Success = false
		returnModel.Token = ""
		returnModel.Msg = conf.UserDisabled
		dblogs_service.Add(conf.LOGIN_FAILE, "USER TRY LOGIN AND USER DISABLED "+user.NickName+" ( From : "+fromLogin+" )", "", ipRequest)
		return returnModel
	}

	tokenGet, err := util.GenerateToken("keep_trying _thief.", "keep_trying _thief.", remember, fromLogin)
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
		if fromLogin == "PHONE_MOVIL" {
			timeLoggout = time.Now().Local().Add(time.Minute * time.Duration(setting.AppSetting.TimeNotPersistToken*100))

		}

	}

	//add session //generate session User
	newSession := authinterfaces.SessionUser{
		Token:          tokenGet,
		Active:         true,
		DateAdd:        time.Now(),
		IdCompany:      user.IdCompany,
		IdUser:         user.ID,
		Remember:       remember,
		TokenExpire:    timeLoggout,
		LastUpdateTime: time.Now(),
	}

	dbsession_user_service.Add(newSession)
	dblogs_service.Add(conf.LOGIN_SUCCESS, "USER LOGIN SUCCESS.. "+user.NickName+" ( From : "+fromLogin+" )", "", ipRequest)
	claim := dbsession_user_service.GetClaimForToken(tokenGet)

	resp := make(map[string]interface{})
	resp["company"] = claim.Company_sys
	resp["user"] = claim.User_sys
	resp["privileges"] = claim.UserPrivileges_sys

	returnModel.Success = true
	returnModel.Token = tokenGet
	returnModel.Msg = conf.UserSuccess

	return returnModel
}

// @Summary Post ClaimUser
// @Produce  json
// @Tags Auth
//
//	//    @ID Authentication
//
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
		appG.Response(http.StatusOK, e.ERROR_AUTH, sendData)
		return
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		sendData.Success = false
		sendData.Msg = "Token Failed..."
		appG.Response(http.StatusOK, e.ERROR_AUTH, sendData)
		return
	}

	dataClaim := dbsession_user_service.GetClaimForToken(token)
	dataModel := make(map[string]interface{})

	dataModel["company"] = dataClaim.Company_sys
	dataModel["user"] = dataClaim.User_sys
	dataModel["dept"] = dataClaim.DeptUser_sys
	dataModel["privilege"] = dataClaim.UserPrivileges_sys
	dataModel["rol"] = dataClaim.RolUser_sys

	dblogs_service.Add(conf.LOGIN_USER_GETCLAIM, " GET_CLAIM SESSION "+dataClaim.User_sys.NickName, dataClaim.User_sys.ID, ipRequest)

	appG.Response(http.StatusOK, e.SUCCESS, dataModel)
}

// @Summary Post ClaimUser
// @Produce  json
// @Tags Auth
//
//	//    @ID Authentication
//
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
		appG.Response(http.StatusOK, e.ERROR_AUTH, sendData)
		return
	}
	//fmt.Println("Here")
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		sendData.Success = false
		sendData.Msg = "Token Failed..."
		appG.Response(http.StatusOK, e.ERROR_AUTH, sendData)
		return
	}

	logout := "false"
	session := dbsession_user_service.FindToToken(token)
	if session.Active {
		session.Active = false
	}
	result := dbsession_user_service.UpdateOne(session)
	if result {
		logout = "true"
	}

	dblogs_service.Add(conf.LOGIN_USER_LOGOUT, "LOGOUT USER "+session.IdUser, session.IdUser, ipRequest)

	myJsonString := `{"logout":` + logout + `}`
	appG.Response(http.StatusOK, e.SUCCESS, myJsonString)
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
		// appG.Response(http.StatusOK, e.ERROR_AUTH, myJsonString)
		// return
		valid = false
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		// appG.Response(http.StatusOK, e.ERROR_AUTH, myJsonString)
		// return
		valid = false
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
	activeToken.LastUpdateTime = time.Now()
	if !activeToken.Active {
		valid = false
	}

	dblogs_service.Add(conf.LOGIN_USER_CHECK_TOKEN, "CHECK TOKEN STATUS: "+strconv.FormatBool(valid), activeToken.IdUser, ipRequest)

	//update lastTimeOnline
	dbsession_user_service.UpdateOne(activeToken)
	appG.Response(http.StatusOK, e.SUCCESS, valid)
}

// @Summary Post Check Password User
// @Produce  json
// @Tags Auth
// @Param username query string true "username"
// @Param password query string true "password"
//
//	//    @ID Authentication
//
// @Security bearerAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth/checkpassword [post]
func PostCheckPasswordUser(c *gin.Context) {

	appG := app.Gin{C: c}
	var valid = false
	if !CheckBearer(c) {
		appG.Response(http.StatusNonAuthoritativeInfo, e.ERROR_AUTH, valid)
		return
	}

	ipRequest = c.ClientIP()

	jsonRequest := models.CheckPasswordUserSys{}
	c.BindJSON(&jsonRequest) //get params from Body

	username := jsonRequest.Username //c.PostForm("username")
	password := jsonRequest.Password //c.Query("password") //c.PostForm("password")

	userFind := dbusers_service.CheckUserPasswordForUser(username, password)
	if userFind.NickName != "" {
		valid = true
	} else {
		userFind = dbusers_service.CheckUserPasswordForEmail(username, password)
		if userFind.NickName != "" {
			valid = true
		}
	}

	//fmt.Println(userFind)
	appG.Response(http.StatusOK, e.SUCCESS, valid)
	//return
}

func CheckBearer(c *gin.Context) bool {

	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		//appG.Response(http.StatusOK, e.ERROR_AUTH, false)
		return false
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		//appG.Response(http.StatusOK, e.ERROR_AUTH, false)
		return false
	}
	return true
}

// @Summary Post CheckIdDevice
// @Produce  json
// @Tags Auth
// @Param idDevice query string true "idDevice"
// @Param type query string true "type"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Router /auth/checkIdDevice [post]
func PostCheckIdDevice(c *gin.Context) {

	appG := app.Gin{C: c}
	valid := validation.Validation{}
	sendResponse := false

	ipRequest = c.ClientIP()
	modelDevice := deviceId{}
	c.BindJSON(&modelDevice) //get params from Body

	modelUser := deviceId{IdDevice: modelDevice.IdDevice, Type: modelDevice.Type}
	ok, _ := valid.Valid(&modelUser)

	fmt.Print("HERE" + strconv.FormatBool(ok))
	fmt.Println(modelUser.Type)
	fmt.Println(modelUser.IdDevice)
	if !ok {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if modelUser.IdDevice == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	found := dbDevicesSys_service.FindToIdDevice(modelDevice.IdDevice)
	if found.Key == "" {
		//is not Device
		newDevice := authinterfaces.DevicesSys_sys{Key: modelDevice.IdDevice,
			IdCompany: "", DateAdd: time.Now(), Status: 0, Type: modelDevice.Type,
			Note: "NUEVO DEVICE BLOCK"}
		dbDevicesSys_service.Add(newDevice)
		dblogs_service.Add(conf.SYSTEM_DEVICE_EVENT, "ADD NEW DEVICE BLOCK TYPE==>  "+modelDevice.Type, modelDevice.IdDevice, ipRequest)
	} else {
		if found.Status == 1 {
			sendResponse = true
		}
		dblogs_service.Add(conf.SYSTEM_DEVICE_EVENT, "ADD CONTACT TO DEVICE IS => "+modelDevice.Type+" result: "+strconv.FormatBool(sendResponse), modelDevice.IdDevice, ipRequest)
	}

	// 	//fmt.Println(result)
	myJsonString := `{"active":` + strconv.FormatBool(sendResponse) + `}`
	fmt.Println(myJsonString)
	appG.Response(http.StatusOK, e.SUCCESS, myJsonString)
	return

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
