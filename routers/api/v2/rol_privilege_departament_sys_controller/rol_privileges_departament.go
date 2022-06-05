package rol_privilege_departament_sys_controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/CastroEduardo/golang-api-rest/models/authinterfaces"
	"github.com/CastroEduardo/golang-api-rest/pkg/app"
	"github.com/CastroEduardo/golang-api-rest/pkg/e"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbdepartamentuser_service"
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

// @Summary Manage Departament Sys
// @Produce  json
// @Tags  Api-v2
// @Param typeOperation query string true "typeOperation"
// @Param idParam query string true "idParam"
// @Param modelJson query string true "modelJson"
// @Security bearerAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/manageDepartament [post]
func ManageDepartamentSys(c *gin.Context) {

	fmt.Println(" === MANAGE DEPARTAMENT ===")
	appG := app.Gin{C: c}
	ipRequest = c.ClientIP()
	//Request Body
	paramRequest := RequestParams{}

	c.BindJSON(&paramRequest)

	fmt.Println(paramRequest)

	auth := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")
	idCompany := dbsession_user_service.GetClaimForToken(token)

	switch paramRequest.TypeOperation {
	case "search":
		//search all
		if paramRequest.IdParam == "" {
			result := dbdepartamentuser_service.GetListFromIdCompany(idCompany.Company.ID)
			u, _ := json.Marshal(result)
			appG.Response(http.StatusOK, e.SUCCESS, string(u))

			//fmt.Println(string(u))
			return
		} else {
			result := dbdepartamentuser_service.FindToId(paramRequest.IdParam)
			u, _ := json.Marshal(result)
			appG.Response(http.StatusOK, e.SUCCESS, string(u))
			return
		}
	case "add":
		//newStr := strings.Replace(paramRequest.ModelData, "'", `"`, -1)
		err := json.Unmarshal([]byte(paramRequest.ModelData), &paramRequest)
		//fmt.Println(paramRequest)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, "ERROR JSON MODEL")
			return
		}
		modelNew := authinterfaces.DepartamentUserSys{}
		json.Unmarshal([]byte(paramRequest.ModelData), &modelNew)

		modelNew.ID = ""
		modelNew.IdCompany = idCompany.Company.ID

		idNew := dbdepartamentuser_service.Add(modelNew)
		//fmt.Println(idNew)

		appG.Response(http.StatusOK, e.SUCCESS, idNew)
		return
	case "update":

		if paramRequest.IdParam == "" {
			appG.Response(http.StatusInternalServerError, e.ERROR, "SEND ID TO UPDATE")
			return
		}

		//newStr := strings.Replace(paramRequest.ModelData, `"`, `'`, -1)
		err := json.Unmarshal([]byte(paramRequest.ModelData), &paramRequest)
		if err != nil {
			fmt.Println("pass")
			appG.Response(http.StatusInternalServerError, e.ERROR, "ERROR JSON MODEL TO UPDATE")
			return
		}

		//fmt.Println(paramRequest.ModelData)
		oldDepart := dbdepartamentuser_service.FindToId(paramRequest.IdParam)
		if oldDepart.IdCompany != idCompany.Company.ID {
			appG.Response(http.StatusInternalServerError, e.ERROR, "FAILD DEPARTAMENT")
			return
		}

		modelNew := authinterfaces.DepartamentUserSys{}
		json.Unmarshal([]byte(paramRequest.ModelData), &modelNew)

		//add updates
		oldDepart.Name = modelNew.Name
		oldDepart.Note = modelNew.Note
		oldDepart.Status = modelNew.Status

		resultUpdate := dbdepartamentuser_service.UpdateOne(oldDepart)
		if resultUpdate {
			appG.Response(http.StatusOK, e.SUCCESS, resultUpdate)
			return
		}

	case "delete":
		if paramRequest.IdParam == "" {
			appG.Response(http.StatusInternalServerError, e.ERROR, "SEND ID TO REMOVE")
			return
		}

		resultRemove := dbdepartamentuser_service.Delete(paramRequest.IdParam)
		if resultRemove {
			appG.Response(http.StatusOK, e.SUCCESS, resultRemove)
			return
		}

	default:
		break
	}

	appG.Response(http.StatusInternalServerError, e.ERROR, "ERROR DEPARTAMENT")
	return
}

type RequestParams struct {
	IdParam       string `json:"idParam"`
	TypeOperation string `json:"typeOperation"`
	ModelData     string `json:"modelData"`
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
