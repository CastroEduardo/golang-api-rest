package users_sys_controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CastroEduardo/golang-api-rest/conf"
	"github.com/CastroEduardo/golang-api-rest/models/authinterfaces"
	"github.com/CastroEduardo/golang-api-rest/pkg/app"
	"github.com/CastroEduardo/golang-api-rest/pkg/e"
	"github.com/CastroEduardo/golang-api-rest/pkg/util"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dblogs_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbsession_user_service"
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
// @Router /api/v2/managedUserSys [post]
func ManagedUserSys(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")
	claimSession := dbsession_user_service.GetClaimForToken(token)

	//fmt.Println(" === ADD USER SYSTEM --")
	appG := app.Gin{C: c}
	ipRequest = c.ClientIP()

	paramRequest := RequestParams{}
	c.BindJSON(&paramRequest)
	switch paramRequest.TypeOperation {
	case "search":

		fmt.Println("SEARCH ___")

		var data []byte
		if paramRequest.IdParam == "" {
			result := dbusers_service.GetListFromIdCompany(claimSession.Company_sys.ID)
			u, _ := json.Marshal(result)
			data = u

		} else {
			result := dbusers_service.GetListFromIdCompany(claimSession.Company_sys.ID)
			u, _ := json.Marshal(result)
			data = u

		}

		//fmt.Println(string(data))
		appG.Response(http.StatusOK, e.SUCCESS, string(data))
		return

	case "add":
		fmt.Println("___ADD___")

		modelRequest := authinterfaces.User_sys{}
		json.Unmarshal([]byte(paramRequest.ModelJson), &modelRequest)

		modelNew := authinterfaces.User_sys{
			NickName:        strings.ToLower(modelRequest.NickName),
			Name:            strings.ToLower(modelRequest.NickName),
			LastName:        "",
			Contact:         "",
			City:            "",
			Gender:          "",
			Email:           strings.ToLower(modelRequest.Email),
			IdDept:          modelRequest.IdDept,
			IdCompany:       claimSession.Company_sys.ID,
			Status:          1,
			Image:           "user.png",
			Note:            modelRequest.Note,
			ForcePass:       true,
			Public:          1,
			Password:        util.Encript([]byte(modelRequest.Password)),
			LastLogin:       time.Now(),
			DefaultPathHome: "/dashboard/workbench",
			DateAdd:         time.Now(),
			ToursInit:       true,
		}

		result := dbusers_service.Add(modelNew)
		fmt.Println(modelNew)
		appG.Response(http.StatusOK, e.SUCCESS, result)
		return
	case "update":

		modelRequest := authinterfaces.User_sys{}
		json.Unmarshal([]byte(paramRequest.ModelJson), &modelRequest)

		userToUpdate := dbusers_service.FindToId(paramRequest.IdParam)

		userToUpdate.Email = strings.ToLower(modelRequest.Email)
		userToUpdate.Note = modelRequest.Note
		userToUpdate.IdDept = modelRequest.IdDept

		result := dbusers_service.UpdateOne(userToUpdate)

		fmt.Println(result)
		appG.Response(http.StatusOK, e.SUCCESS, result)
		return
	case "delete":

		fmt.Println("___DELeTE___")
		userFind := dbusers_service.FindToId(paramRequest.IdParam)
		out, err := json.Marshal(userFind)
		if err != nil {
			panic(err)
		}
		response := dbusers_service.DeleteToId(paramRequest.IdParam)
		go dblogs_service.Add(conf.LOGIN_USER_EVENT_REMOVE, "DELETE USER: "+string(out), claimSession.User_sys.ID, ipRequest)
		appG.Response(http.StatusOK, e.SUCCESS, response)
		return
	case "remove-tours-init":

		fmt.Println("DISaBLE_TOURS")
		idUser := claimSession.User_sys.ID
		userUpdate := dbusers_service.FindToId(idUser)
		userUpdate.ToursInit = false
		result := dbusers_service.UpdateOne(userUpdate)
		appG.Response(http.StatusOK, e.SUCCESS, result)
		return
	case "isAccount":
		isAccount := dbusers_service.IsAccount(paramRequest.IdParam)
		appG.Response(http.StatusOK, e.SUCCESS, isAccount)
		go dblogs_service.Add(conf.USER_CHECK_ISACCOUNT, "CHECK ACCOUNT EXIST RESULT :  "+strconv.FormatBool(isAccount)+" =>> "+paramRequest.IdParam, claimSession.User_sys.ID, ipRequest)
		return
	default:
		break
	}

	appG.Response(http.StatusInternalServerError, e.ERROR, "false")
	return

}

type RequestParams struct {
	IdParam       string `json:"idParam"`
	TypeOperation string `json:"typeOperation"`
	ModelJson     string `json:"modelJson"`
}
