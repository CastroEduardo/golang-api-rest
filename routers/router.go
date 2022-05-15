package routers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CastroEduardo/golang-api-rest/conf"
	_ "github.com/CastroEduardo/golang-api-rest/docs"
	"github.com/CastroEduardo/golang-api-rest/middleware/jwt"
	"github.com/CastroEduardo/golang-api-rest/routers/api"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/CastroEduardo/golang-api-rest/pkg/export"
	"github.com/CastroEduardo/golang-api-rest/pkg/qrcode"
	"github.com/CastroEduardo/golang-api-rest/pkg/upload"
	v1 "github.com/CastroEduardo/golang-api-rest/routers/api/v1"
	v2 "github.com/CastroEduardo/golang-api-rest/routers/api/v2"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(CORSMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	fmt.Println("-- LOADING --")

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	//r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	// staticDir := "/upload/images"
	// r.Handle(staticDir, http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))

	r.POST("/auth", api.PostAuth)
	r.GET("/auth", api.GetAuth)
	r.POST("/auth/claim-user", api.PostClaimUser)
	r.POST("/auth/logout", api.Postlogout)

	r.POST("/auth/checkstatustoken", api.PostCheckStatusSession)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// r.POST("/upload", api.UploadImage)

	apiv1 := r.Group(conf.NameUrlApi1)
	apiv1.Use(jwt.JWT())
	{
		apiv1.GET(conf.ArticlesParms_GET, v1.GetArticle)
	}

	apiv2 := r.Group(conf.NameUrlApi2)
	apiv2.Use(jwt.JWT())
	{
		//src := upload.GetImageFullPath()
		//apiv1.Static("/upload/images", src)
		//apiv2.StaticFS("/test/", http.Dir(src))
		src := upload.GetImageFullPath()
		fmt.Println(src)
		apiv2.StaticFS("/upload/images", http.Dir(src))
		//apiv2.Static("/upload/images", "/Users/ecastro_mbp/Desktop/golang_restapi/upload/images/")
		apiv2.GET(conf.ArticlesParms_GET, v2.GetArticle)

	}
	// 	//获取标签列表
	// 	apiv1.GET("/tags", v1.GetTags)
	// 	//新建标签
	// 	apiv1.POST("/tags", v1.AddTag)
	// 	//更新指定标签
	// 	apiv1.PUT("/tags/:id", v1.EditTag)
	// 	//删除指定标签
	// 	apiv1.DELETE("/tags/:id", v1.DeleteTag)
	// 	//导出标签
	// 	r.POST("/tags/export", v1.ExportTag)
	// 	//导入标签
	// 	r.POST("/tags/import", v1.ImportTag)

	// 	//获取文章列表
	// 	apiv1.GET("/articles", v1.GetArticles)
	// 	//获取指定文章
	// 	apiv1.GET("/articles/:id", v1.GetArticle)
	// 	//新建文章
	// 	apiv1.POST("/articles", v1.AddArticle)
	// 	//更新指定文章
	// 	apiv1.PUT("/articles/:id", v1.EditArticle)
	// 	//删除指定文章
	// 	apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	// 	//生成文章海报
	// 	apiv1.POST("/articles/poster/generate", v1.GenerateArticlePoster)
	// }
	// apiv2 := r.Group("/api/v2")
	// apiv2.Use(jwt.JWT())
	// {
	// 	//获取标签列表
	// 	apiv2.GET("/tags", v2.GetTags)
	// 	//新建标签

	// }

	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
