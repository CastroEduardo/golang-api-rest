package uploadsfiles

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/CastroEduardo/golang-api-rest/conf"
	"github.com/CastroEduardo/golang-api-rest/pkg/app"
	"github.com/CastroEduardo/golang-api-rest/pkg/e"
	"github.com/CastroEduardo/golang-api-rest/pkg/upload"
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

// @Summary Upload Files to Server
// @Produce  json
// @Tags  Api-v2
// @Param modelUser query string true "modelUser"
// @Security bearerAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/manageduploadfiles [post]
func ManagedUploads(c *gin.Context) {

	auth := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")
	claimSession := dbsession_user_service.GetClaimForToken(token)

	appG := app.Gin{C: c}
	ipRequest = c.ClientIP()
	src := upload.GetImageFullPath()

	var newNameFile = "" //util.GetUniqueId()
	var typeAction = ""
	var idAction = ""
	var typeFolder = ""

	var pathNew = ""
	// Multipart form
	form, _ := c.MultipartForm()
	files := form.File["file"]

	for _, file := range files {
		divideName := strings.Split(file.Filename, "*")
		file.Filename = divideName[0]
		typeAction = divideName[1]
		idAction = divideName[2]
		fmt.Println(file.Filename)
		fmt.Println(typeAction)
		fmt.Println(idAction)

		fileContent, _ := file.Open()
		var byteContainer []byte
		byteContainer = make([]byte, 1000000)
		fileContent.Read(byteContainer)
		// fmt.Println(byteContainer)

		//var asd:=interface{}
		//form-data; name="file"; filename="avatar.png"
		//fmt.Println(file.Header[0])
		newNameFile = util.GetUniqueId() + upload.GetImageName(file.Filename)
		typeFile := strings.Split(file.Filename, ".")
		if typeFile[1] == "png" || typeFile[1] == "jpg" || typeFile[1] == "jpeg" || typeFile[1] == "svg" {
			typeFolder = "/images/"
			pathNew = src + claimSession.Company_sys.FolderFiles + typeFolder + newNameFile

			//filename, err :=  u  (buffer, 40, "uploads")

		} else {
			//dst := path.Join(src, file.Filename)
			// Upload the file to specific dst.
			typeFolder = "/files/"
			pathNew = src + claimSession.Company_sys.FolderFiles + typeFolder + newNameFile
		}
		error := c.SaveUploadedFile(file, pathNew)
		if error != nil {
			newNameFile = ""
			appG.Response(http.StatusInternalServerError, e.ERROR, newNameFile)
			return
		}

		//Actions to File of Case
		var oldNameFileToRemove = ""
		switch typeAction {
		case "updateAvatar":
			userToUpdate := dbusers_service.FindToId(claimSession.User_sys.ID)
			oldNameFileToRemove = userToUpdate.Image
			userToUpdate.Image = newNameFile
			result := dbusers_service.UpdateOne(userToUpdate)
			if result {
				go dblogs_service.Add(conf.USER_EVENT_UPDATE_AVATAR, "USER UPDATE AVATAR  : "+claimSession.User_sys.NickName, claimSession.User_sys.ID, ipRequest)
			}

			if oldNameFileToRemove != "user.png" {
				oldPath := src + claimSession.Company_sys.FolderFiles + typeFolder + oldNameFileToRemove
				newPath := src + claimSession.Company_sys.FolderFiles + "/remove/" + oldNameFileToRemove
				go util.MoveFile(oldPath, newPath)
				go dblogs_service.Add(conf.SYSTEM_EVENT_MOVE_FILE, "SYSTEM_MOVE_FILE "+"*"+oldNameFileToRemove+"*"+" AFTER TO UPDATE AVATAR : "+claimSession.User_sys.NickName, claimSession.User_sys.ID, ipRequest)
			}

		default:

		}

		//fmt.Println(pathNew)
		appG.Response(http.StatusOK, e.SUCCESS, newNameFile)
		return
	}

	//appG.Response(http.StatusOK, e.SUCCESS, "false")

}
