package routers

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-crawler/conf"
	service "go-crawler/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(sig ...os.Signal) {
	router := SetupRouter()

	if len(sig) == 0 {
		sig = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	signalChan := make(chan os.Signal, 1)

	go func() {
		router.Run(fmt.Sprintf(":%v", conf.GetEnv("HTTP_PORT")))
	}()
	signal.Notify(signalChan, sig...)
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.ExposeHeaders = []string{"Authorization"}
	config.AllowCredentials = true
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	router.GET("/ping", service.PingHandler)
	router.GET("/task", service.TaskHandler)
	router.POST("/query_blogs", service.QueryBlogsHandler)
	router.POST("/add_bloggers", service.AddBloggersHandler)
	category := router.Group("/category")
	{
		category.POST("/set", service.SetCategoriesHandler)
		category.GET("/query", service.QueryCategroyHandler)
		category.POST("/set_name", service.SetCategoriesNameHandler)
	}
	tags := router.Group("/tags")
	{
		tags.GET("/query", service.QueryTagsHandler)
		//tags.POST("/cache", service.SaveToRedisHandler)
		tags.POST("/cache_keywords", service.SaveKeywordToRedisHandler)
		tags.POST("/keywords", service.QueryKeywordsHandler)
		tags.POST("/set_keywords", service.SetKeywordsHandler)
	}
	return router
}
