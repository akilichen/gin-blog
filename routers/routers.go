package routers

import (
	"gin-blog/middleware/jwt"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/upload"
	"gin-blog/routers/api"
	"gin-blog/routers/api/v1"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	r.GET("/auth", api.GetAuth) //获取accessToken

	apiOfV1 := r.Group("/api/v1")
	apiOfV1.Use(jwt.JWT())
	{
		//获取标签列表
		apiOfV1.GET("/tags", v1.GetTags)
		//新建标签
		apiOfV1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiOfV1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiOfV1.DELETE("/tags/:id", v1.DeleteTag)

		//获取文章列表
		apiOfV1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiOfV1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiOfV1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiOfV1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiOfV1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload", api.UploadImage)

	//使用golang自建静态文件服务器，一般在实战中都是采用CDN或者分布式文件系统
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

	return r
}
