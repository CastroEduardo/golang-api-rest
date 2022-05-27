package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	//github.com/CastroEduardo/golang-api-rest/models"

	// "github.com/CastroEduardo/golang-api-rest/pkg/gredis"

	"github.com/CastroEduardo/golang-api-rest/pkg/logging"
	"github.com/CastroEduardo/golang-api-rest/pkg/mongo_db"
	"github.com/CastroEduardo/golang-api-rest/pkg/setting"
	"github.com/CastroEduardo/golang-api-rest/pkg/util"
	"github.com/CastroEduardo/golang-api-rest/routers"
)

func init() {
	setting.Setup()
	//models.Setup()
	logging.Setup()
	// gredis.Setup()
	mongo_db.Setup()
	util.Setup()

	time.Sleep(1 * time.Second)

	//setting_create_user.Create_first_user()
	//result := logs_service.Add("ads")

}

// @title Golang Gin API
// @version 1.0
// @description Sample of WebApiRest
// @termsOfService https://github.com/CastroEduardo/golang-api-rest
// @license.name MIT
// @license.url https://github.com/CastroEduardo/golang-api-rest/blob/master/LICENSE
// @securityDefinitions.apikey bearerAuth
// @in header
// @name Authorization
func main() {
	gin.SetMode(setting.ServerSetting.RunMode)

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf("%s:%d", setting.ServerSetting.Host, setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	fmt.Println(endPoint)
	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
		// 	TLSConfig:      cfg,
		// 	TLSNextProto:   make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	//server.ListenAndServe()
	//finally we Listen for requests and serve them up on a specific port

	server.ListenAndServeTLS("192.168.0.101+3.pem", "192.168.0.101+3-key.pem")

	//server.ListenAndServe()
	// srv := &http.Server{
	// 	Addr:         "",
	// 	Handler:      routersInit,
	// 	TLSConfig:    cfg,
	// 	TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	// }

	//server.TLSConfig()

	// server.ListenAndServe()

	//server.ListenAndServeTLS("", "")

	// If you want Graceful Restart, you need a Unix system and download github.com/fvbock/endless
	//endless.DefaultReadTimeOut = readTimeout
	//endless.DefaultWriteTimeOut = writeTimeout
	//endless.DefaultMaxHeaderBytes = maxHeaderBytes
	//server := endless.NewServer(endPoint, routersInit)
	//server.BeforeBegin = func(add string) {
	//	log.Printf("Actual pid is %d", syscall.Getpid())
	//}
	//
	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Printf("Server err: %v", err)
	//}
}
