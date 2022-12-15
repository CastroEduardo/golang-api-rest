package option_sys_controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"strings"

	"github.com/CastroEduardo/golang-api-rest/conf"
	"github.com/CastroEduardo/golang-api-rest/models/authinterfaces"
	"github.com/CastroEduardo/golang-api-rest/models/interfaces_public"
	"github.com/CastroEduardo/golang-api-rest/pkg/app"
	"github.com/CastroEduardo/golang-api-rest/pkg/e"
	"github.com/CastroEduardo/golang-api-rest/pkg/util"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbgeneric_option_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbsession_user_service"
	"github.com/gin-gonic/gin"
	//"github.com/boombuler/barcode/qr"
	//github.com/CastroEduardo/golang-api-rest/pkg/qrcode"
	//github.com/CastroEduardo/golang-api-rest/pkg/setting"
	//github.com/CastroEduardo/golang-api-rest/pkg/util"
	//github.com/CastroEduardo/golang-api-rest/service/tag_service"
)

var ipRequest = ""

// @Summary ManagedGenericOptions Sys
// @Produce  json
// @Tags  Api-v2
// @Param modelOptions query string true "modelOptions"
// @Security bearerAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/managedGenericOptions [post]
func ManagedGeneryOptions(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")
	claimSession := dbsession_user_service.GetClaimForToken(token)

	//fmt.Println(" === ADD USER SYSTEM --")
	appG := app.Gin{C: c}
	ipRequest = c.ClientIP()
	fmt.Println("HERE")
	paramRequest := RequestParams{}
	c.BindJSON(&paramRequest)

	fmt.Println(paramRequest.Type)

	switch paramRequest.Type {
	case conf.GENERIC_FARM:

		farm := managedFarm(paramRequest, claimSession)
		appG.Response(http.StatusOK, e.SUCCESS, farm)
		return
	default:
		allItems := managedAll(paramRequest, claimSession)
		appG.Response(http.StatusOK, e.SUCCESS, allItems)
		return
	}

	// appG.Response(http.StatusInternalServerError, e.ERROR, "false")
	// return

}

func managedFarm(params RequestParams, claim_Session authinterfaces.ClaimSession) string {
	switch params.TypeOperation {
	case "add":
		fmt.Println("__ADD_")
		modelRequest := interfaces_public.GenericList{}

		json.Unmarshal([]byte(params.ModelJson), &modelRequest)
		fmt.Println(modelRequest.Name)

		newFarm := interfaces_public.GenericList{
			Name:      modelRequest.Name,
			IdKey:     util.GetUniqueId(),
			Identity:  conf.GENERIC_FARM,
			Status:    modelRequest.Status,
			Value1:    modelRequest.Value1,
			Value2:    modelRequest.Value2,
			Value3:    modelRequest.Value3,
			Value4:    modelRequest.Value4,
			Value5:    modelRequest.Value5,
			Value6:    modelRequest.Value6,
			Note:      modelRequest.Note,
			Date:      time.Now(),
			IdCompany: claim_Session.User_sys.IdCompany,
		}
		return dbgeneric_option_service.Add(newFarm)

	case "update":

		item := interfaces_public.GenericList{}
		json.Unmarshal([]byte(params.ModelJson), &item)
		//fmt.Println(item.Name)
		//item.Name = "CAMbio"
		update := dbgeneric_option_service.UpdateOne(item)

		return strconv.FormatBool(update)

	case "search":
		fmt.Println("_Search_")
		var data []byte
		if params.IdParam == "#*#" {
			//fmt.Println(claim_Session.User_sys.IdCompany)
			result := dbgeneric_option_service.GetListFromIdCompany(claim_Session.User_sys.IdCompany)
			data, _ = json.Marshal(result)
			//fmt.Println(string(data))
			return string(data)
		} else {

		}
		break

	default:
		break

	}

	return ""

}

// ! All only
func managedAll(params RequestParams, claim_Session authinterfaces.ClaimSession) string {
	switch params.TypeOperation {
	case "add":

	case "search":
		fmt.Println("_Search_All_")
		var data []byte
		if params.IdParam == "#*#" {
			//fmt.Println("ALL")
			// fmt.Println(claim_Session.User_sys.IdCompany)
			result := dbgeneric_option_service.GetListFromIdCompany(claim_Session.User_sys.IdCompany)
			data, _ = json.Marshal(result)
			//fmt.Println(string(data))
			return string(data)
		} else {

		}
		break

	default:
		break

	}

	return ""

}

type RequestParams struct {
	Type          string `json:"type"`
	IdParam       string `json:"idParam"`
	TypeOperation string `json:"typeOperation"`
	ModelJson     string `json:"modelJson"`
}
