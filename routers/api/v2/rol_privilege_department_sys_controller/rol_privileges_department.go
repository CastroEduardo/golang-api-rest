package rol_privilege_department_sys_controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/CastroEduardo/golang-api-rest/conf"
	"github.com/CastroEduardo/golang-api-rest/models/authinterfaces"
	"github.com/CastroEduardo/golang-api-rest/pkg/app"
	"github.com/CastroEduardo/golang-api-rest/pkg/e"
	"github.com/CastroEduardo/golang-api-rest/pkg/util"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbPrivilegeuser_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbdepartmentuser_service"
	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dblogs_service"

	"github.com/CastroEduardo/golang-api-rest/service/mongo_service/dbrol_user_service"
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

// @Summary Manage DepartamentUser Sys
// @Produce  json
// @Tags  Api-v2
// @Param typeOperation query string true "typeOperation"
// @Param idParam query string true "idParam"
// @Param modelJson query string true "modelJson"
// @Security bearerAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/manageDepartment [post]
func ManageDepartamentSys(c *gin.Context) {

	appG := app.Gin{C: c}
	ipRequest = c.ClientIP()
	// //Request Body
	paramRequest := RequestParams{}
	c.BindJSON(&paramRequest)

	auth := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")
	claimSession := dbsession_user_service.GetClaimForToken(token)
	newList := []authinterfaces.DptsUser_sys{}

	switch paramRequest.TypeOperation {
	case "search":
		//search all
		if paramRequest.IdParam == "" {
			result := dbdepartmentuser_service.GetListFromIdCompany(claimSession.Company_sys.ID)
			u, _ := json.Marshal(result)
			appG.Response(http.StatusOK, e.SUCCESS, string(u))
			// fmt.Println(result)

		} else {
			result := dbdepartmentuser_service.FindToId(paramRequest.IdParam)
			u, _ := json.Marshal(result)
			appG.Response(http.StatusOK, e.SUCCESS, string(u))
		}

		go dblogs_service.Add(conf.LOGIN_USER_SEARCH, "SEARCH DEPT SYS", claimSession.User_sys.ID, ipRequest)
		return
	case "add":
		//newStr := strings.Replace(paramRequest.ModelJson, "'", `"`, -1)
		err := json.Unmarshal([]byte(paramRequest.ModelJson), &paramRequest)
		//fmt.Println(paramRequest)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, "ERROR JSON MODEL")
			return
		}

		modelNew := authinterfaces.DptsUser_sys{}
		json.Unmarshal([]byte(paramRequest.ModelJson), &modelNew)

		if paramRequest.IdParam == "" {
			//is Parent
			newDpt := authinterfaces.DptsUser_sys{
				Title:     modelNew.Title,
				Key:       "key",
				Note:      modelNew.Note,
				Status:    modelNew.Status,
				OrderNo:   1,
				IdRole:    modelNew.IdRole,
				IdCompany: claimSession.User_sys.IdCompany,
				DateAdd:   time.Now(),
				Children:  []authinterfaces.DptsUser_sys{},
			}
			idParent := dbdepartmentuser_service.Add(newDpt)
			getDptUpdate := dbdepartmentuser_service.FindToId(idParent)
			getDptUpdate.Key = idParent
			dbdepartmentuser_service.UpdateOne(getDptUpdate)

		} else {
			//is Children
			var chilAdd []authinterfaces.DptsUser_sys
			parentId, childId := dbdepartmentuser_service.SearchParentIdCompany(claimSession.Company_sys.ID, paramRequest.IdParam)
			if childId != "" {
				//ADD INTO TO CHILD
			} else {
				//ADD INTO TO PARENT
				getDptUpdate1 := dbdepartmentuser_service.FindToKey(parentId)
				newDptChilder := authinterfaces.DptsUser_sys{
					Title:     modelNew.Title,
					Key:       util.GetUniqueId(),
					Note:      modelNew.Note,
					Status:    modelNew.Status,
					OrderNo:   1,
					IdRole:    modelNew.IdRole,
					IdCompany: claimSession.User_sys.IdCompany,
					DateAdd:   time.Now(),
					Children:  nil,
				}
				chilAdd = append(getDptUpdate1.Children, newDptChilder)
				getDptUpdate1.Children = chilAdd
				dbdepartmentuser_service.UpdateOneToKey(getDptUpdate1)
				//fmt.Println(updtResult)
			}
		}

		go dblogs_service.Add(conf.LOGIN_USER_EVENT_ADD, "ADD : D P T ", claimSession.User_sys.ID, ipRequest)

		appG.Response(http.StatusOK, e.SUCCESS, "")
		return
	case "update":
		if paramRequest.IdParam == "" {
			appG.Response(http.StatusInternalServerError, e.ERROR, "SEND ID TO UPDATE")
			return
		}

		//newStr := strings.Replace(paramRequest.ModelJson, `"`, `'`, -1)
		err := json.Unmarshal([]byte(paramRequest.ModelJson), &paramRequest)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, "ERROR JSON MODEL TO UPDATE")
			return
		}

		idParent, idChild := dbdepartmentuser_service.SearchParentIdCompany(claimSession.Company_sys.ID, paramRequest.IdParam)

		modelNew := authinterfaces.DptsUser_sys{}
		json.Unmarshal([]byte(paramRequest.ModelJson), &modelNew)

		if idChild != "" {
			//is CHILD
			//fmt.Println("IS CHILD", idParent)
			parentToupdate := dbdepartmentuser_service.FindToKey(idParent)
			for _, v := range parentToupdate.Children {
				if idChild == v.Key {
					//fmt.Println(" ****1 " + v.Key)
					v.Note = modelNew.Note
					v.Status = modelNew.Status
					v.Title = modelNew.Title
					v.IdRole = modelNew.IdRole
				}

				newList = append(newList, v)
			}
			parentToupdate.Children = newList
			dbdepartmentuser_service.UpdateOneToKey(parentToupdate)
		} else {
			//is PARENT
			toUpdate := dbdepartmentuser_service.FindToKey(paramRequest.IdParam)

			toUpdate.Note = modelNew.Note
			toUpdate.Title = modelNew.Title
			toUpdate.Status = modelNew.Status
			toUpdate.IdRole = modelNew.IdRole
			dbdepartmentuser_service.UpdateOneToKey(toUpdate)
		}
		go dblogs_service.Add(conf.LOGIN_USER_EVENT_UPDATE, "UPDATE : D P T ", claimSession.User_sys.ID, ipRequest)

		appG.Response(http.StatusOK, e.SUCCESS, "true")
		return
	case "delete":
		if paramRequest.IdParam == "" {
			appG.Response(http.StatusInternalServerError, e.ERROR, "SEND ID TO REMOVE")
			return
		}

		idParent, idChild := dbdepartmentuser_service.SearchParentIdCompany(claimSession.Company_sys.ID, paramRequest.IdParam)
		if idChild != "" {
			//is child
			parentToupdate := dbdepartmentuser_service.FindToKey(idParent)
			for _, v := range parentToupdate.Children {
				if idChild != v.Key {
					newList = append(newList, v)
				}
			}
			parentToupdate.Children = newList
			dbdepartmentuser_service.UpdateOneToKey(parentToupdate)

		} else {
			//is Parent
			dbdepartmentuser_service.DeleteToKey(idParent)
		}
		go dblogs_service.Add(conf.LOGIN_USER_EVENT_REMOVE, "DELETE : D P T ", claimSession.User_sys.ID, ipRequest)

		appG.Response(http.StatusOK, e.SUCCESS, "true")
		return
	default:
		break
	}
}

// @Summary Manage RolesUser Sys
// @Produce  json
// @Tags  Api-v2
// @Param typeOperation query string true "typeOperation"
// @Param idParam query string true "idParam"
// @Param modelJson query string true "modelJson"
// @Security bearerAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/managerolsys [post]
func ManageRolSys(c *gin.Context) {

	appG := app.Gin{C: c}
	ipRequest = c.ClientIP()
	//Request Body
	paramRequest := RequestParams{}
	c.BindJSON(&paramRequest)

	auth := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")
	claimSession := dbsession_user_service.GetClaimForToken(token)

	switch paramRequest.TypeOperation {
	case "search":
		//search all
		if paramRequest.IdParam == "" {
			result := dbrol_user_service.GetListFromIdCompany(claimSession.Company_sys.ID)
			u, _ := json.Marshal(result)
			appG.Response(http.StatusOK, e.SUCCESS, string(u))
			//fmt.Println(string(u))

		} else {
			result := dbdepartmentuser_service.FindToId(paramRequest.IdParam)
			u, _ := json.Marshal(result)
			appG.Response(http.StatusOK, e.SUCCESS, string(u))

		}
		go dblogs_service.Add(conf.LOGIN_USER_SEARCH, "SEARCH ROLE SYS", claimSession.User_sys.ID, ipRequest)
		return
	case "add":
		//newStr := strings.Replace(paramRequest.ModelJson, "'", `"`, -1)
		err := json.Unmarshal([]byte(paramRequest.ModelJson), &paramRequest)
		//fmt.Println(paramRequest)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, "ERROR JSON MODEL")
			return
		}
		modelNew := authinterfaces.RolUser_sys{}
		json.Unmarshal([]byte(paramRequest.ModelJson), &modelNew)

		modelNew.ID = ""

		modelNew.IdCompany = claimSession.Company_sys.ID

		idNew := dbrol_user_service.Add(modelNew)
		//fmt.Println(idNew)

		dblogs_service.Add(conf.LOGIN_USER_EVENT_ADD, "ADD : ROLE ", claimSession.User_sys.ID, ipRequest)

		appG.Response(http.StatusOK, e.SUCCESS, idNew)
		return
	case "update":

		if paramRequest.IdParam == "" {
			appG.Response(http.StatusInternalServerError, e.ERROR, "SEND ID TO UPDATE")
			return
		}

		//newStr := strings.Replace(paramRequest.ModelJson, `"`, `'`, -1)
		err := json.Unmarshal([]byte(paramRequest.ModelJson), &paramRequest)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, "ERROR JSON MODEL TO UPDATE")
			return
		}

		//fmt.Println(paramRequest.ModelJson)
		oldRol := dbrol_user_service.FindToId(paramRequest.IdParam)
		if oldRol.IdCompany != claimSession.Company_sys.ID {
			appG.Response(http.StatusInternalServerError, e.ERROR, "FAILD DEPARTAMENT")
			return
		}

		modelNew := authinterfaces.RolUser_sys{}
		json.Unmarshal([]byte(paramRequest.ModelJson), &modelNew)

		//add updates
		oldRol.Name = modelNew.Name
		oldRol.Note = modelNew.Note
		oldRol.Status = modelNew.Status
		oldRol.IdPrivilege = modelNew.IdPrivilege

		resultUpdate := dbrol_user_service.UpdateOne(oldRol)

		dblogs_service.Add(conf.LOGIN_USER_EVENT_UPDATE, "UPDATE : ROLE ", claimSession.User_sys.ID, ipRequest)

		if resultUpdate {
			appG.Response(http.StatusOK, e.SUCCESS, resultUpdate)
			return
		}
	case "delete":
		if paramRequest.IdParam == "" {
			appG.Response(http.StatusInternalServerError, e.ERROR, "SEND ID TO REMOVE")
			return
		}

		resultRemove := dbrol_user_service.Delete(paramRequest.IdParam)

		dblogs_service.Add(conf.LOGIN_USER_EVENT_REMOVE, "REMOVE : ROLE ", claimSession.User_sys.ID, ipRequest)

		if resultRemove {
			appG.Response(http.StatusOK, e.SUCCESS, resultRemove)
			return
		}
	default:
		break
	}

	appG.Response(http.StatusInternalServerError, e.ERROR, "ERROR DEPARTAMENT")

}

// @Summary Manage Privilege Sys
// @Produce  json
// @Tags  Api-v2
// @Param typeOperation query string true "typeOperation"
// @Param idParam query string true "idParam"
// @Param modelJson query string true "modelJson"
// @Security bearerAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/manageprivilegesys [post]
func ManagePrivilegeSys(c *gin.Context) {

	appG := app.Gin{C: c}
	ipRequest = c.ClientIP()
	// //Request Body
	paramRequest := RequestParams{}
	c.BindJSON(&paramRequest)

	auth := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")
	claimSession := dbsession_user_service.GetClaimForToken(token)

	switch paramRequest.TypeOperation {
	case "search":

		//search all
		if paramRequest.IdParam == "" {
			result := dbPrivilegeuser_service.GetListFromIdCompany(claimSession.Company_sys.ID)
			u, _ := json.Marshal(result)
			appG.Response(http.StatusOK, e.SUCCESS, string(u))

			// fmt.Println(result)

		} else {
			result := dbPrivilegeuser_service.FindToId(paramRequest.IdParam)
			u, _ := json.Marshal(result)
			appG.Response(http.StatusOK, e.SUCCESS, string(u))

		}

		go dblogs_service.Add(conf.LOGIN_USER_SEARCH, "SEARCH PRIVILEGE", claimSession.User_sys.ID, ipRequest)
		return

	case "add":

		modelNew := authinterfaces.UserPrivileges_sys{}
		json.Unmarshal([]byte(paramRequest.ModelJson), &modelNew)

		newAdd := authinterfaces.UserPrivileges_sys{
			Name:          modelNew.Name,
			Status:        1,
			IdCompany:     claimSession.Company_sys.ID,
			DefaultPage:   modelNew.DefaultPage,
			ListPages:     modelNew.ListPages,
			ListFunctions: modelNew.ListFunctions,
			DateAdd:       time.Now(),
			Note:          modelNew.Note,
		}
		result := dbPrivilegeuser_service.Add(newAdd)

		dblogs_service.Add(conf.LOGIN_USER_EVENT_ADD, "ADD : PRIVILEGE ", claimSession.User_sys.ID, ipRequest)

		appG.Response(http.StatusOK, e.SUCCESS, result)
		return
	case "update":

		modelNew := authinterfaces.UserPrivileges_sys{}
		json.Unmarshal([]byte(paramRequest.ModelJson), &modelNew)

		fmt.Println(modelNew.ListFunctions)

		itemToUpdate := dbPrivilegeuser_service.FindToId(paramRequest.IdParam)

		itemToUpdate.DefaultPage = modelNew.DefaultPage
		itemToUpdate.Name = modelNew.Name
		itemToUpdate.ListFunctions = modelNew.ListFunctions
		itemToUpdate.ListPages = modelNew.ListPages
		itemToUpdate.Note = modelNew.Note

		dbPrivilegeuser_service.UpdateOne(itemToUpdate)

		dblogs_service.Add(conf.LOGIN_USER_EVENT_UPDATE, "UPDATE : PRIVILEGE ", claimSession.User_sys.ID, ipRequest)

		appG.Response(http.StatusOK, e.SUCCESS, "true")
		return
	case "delete":
		result := dbPrivilegeuser_service.Delete(paramRequest.IdParam)

		dblogs_service.Add(conf.LOGIN_USER_EVENT_REMOVE, "REMOVE : PRIVILEGE ", claimSession.User_sys.ID, ipRequest)

		appG.Response(http.StatusOK, e.SUCCESS, result)
		return

	default:
		break
	}
}

type RequestParams struct {
	IdParam       string `json:"idParam"`
	TypeOperation string `json:"typeOperation"`
	ModelJson     string `json:"modelJson"`
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

	appG := app.Gin{C: c}
	ipRequest = c.ClientIP()
	dbusers_service.FindToId("")
	//fmt.Println(ipRequest)
	appG.Response(http.StatusOK, e.SUCCESS, "Data")
	return
}
