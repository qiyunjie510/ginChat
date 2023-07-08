package router

import (
	"ginchat/controller"
	"ginchat/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	// swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 静态资源
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("view/**/*")

	// 页面
	r.GET("/", controller.GetIndex)
	r.GET("/index", controller.GetIndex)
	r.GET("/register", controller.GetRegister)
	r.GET("/toChat", controller.ToChat)

	// 用户模块
	r.GET("/api/user/list", controller.List)
	r.POST("/api/user/login", controller.Login)
	r.POST("/api/user/register", controller.Register)
	r.POST("/api/user/update", controller.UpdateUser)
	r.GET("/api/user/delete", controller.DeleteUser)

	//
	r.POST("/api/searchFriends", controller.SearchFriends)

	// websocket
	r.GET("/ws/user/chat", controller.Chat)
	r.GET("/ws/user/chatV2", controller.ChatV2)
	return r
}
