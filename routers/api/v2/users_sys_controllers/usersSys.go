package users_sys_controllers

import (
	"fmt"
	"net/http"

	"github.com/CastroEduardo/golang-api-rest/pkg/app"
	"github.com/CastroEduardo/golang-api-rest/pkg/e"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbusers_service"
	"github.com/gin-gonic/gin"
	//"github.com/boombuler/barcode/qr"
	//github.com/CastroEduardo/golang-api-rest/pkg/qrcode"
	//github.com/CastroEduardo/golang-api-rest/pkg/setting"
	//github.com/CastroEduardo/golang-api-rest/pkg/util"
	//github.com/CastroEduardo/golang-api-rest/service/tag_service"
)

var ipRequest = ""

// @Summary Add user to System
// @Produce  json
// @Tags  Api-v2
// @Param modelUser query string true "modelUser"
// @Security bearerAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/addusersys [post]
func AddUserSys(c *gin.Context) {

	fmt.Println(" === ADD USER SYSTEM --")
	appG := app.Gin{C: c}
	ipRequest = c.ClientIP()

	fmt.Println(ipRequest)

	// jsonParsed, err := gabs.ParseJSON(formModel)
	// if err != nil {
	// 	panic(err)
	// }

	appG.Response(http.StatusOK, e.SUCCESS, "DATA")
	return
}

// @Summary Get list userSys
// @Produce  json
// @Tags  Api-v2
// @Param id query string true "id"
// @Security bearerAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/userlistsys [post]
func UserSysList(c *gin.Context) {

	fmt.Println(" === GET LIST USERS SYSTEM --")
	appG := app.Gin{C: c}
	ipRequest = c.ClientIP()

	dbusers_service.FindToId("")

	fmt.Println(ipRequest)

	appG.Response(http.StatusOK, e.SUCCESS, "Data")
	return
}
