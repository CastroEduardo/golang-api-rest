package v2

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/CastroEduardo/golang-api-rest/pkg/app"
	"github.com/CastroEduardo/golang-api-rest/pkg/e"
	"github.com/gin-gonic/gin"
	//"github.com/boombuler/barcode/qr"
	//github.com/CastroEduardo/golang-api-rest/pkg/qrcode"
	//github.com/CastroEduardo/golang-api-rest/pkg/setting"
	//github.com/CastroEduardo/golang-api-rest/pkg/util"
	//github.com/CastroEduardo/golang-api-rest/service/tag_service"
)

var ipRequest = ""

// @Summary Get list userssys 1
// @Produce  json
// @Tags  Api-v2
// @Security bearerAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v2/userlistsys [post]
func UserSysList(c *gin.Context) {

	fmt.Println(" ===HERE --")
	appG := app.Gin{C: c}
	ipRequest = c.ClientIP()

	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		return
	}

	//fmt.Println("Here")
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		return
	}

	fmt.Println(token)

	appG.Response(http.StatusOK, e.SUCCESS, "")
	return

}
