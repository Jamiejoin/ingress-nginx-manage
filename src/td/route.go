package main

import (
	"github.com/gin-gonic/gin"
	. "td/apis"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/gin-contrib/cors"
)


func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/search/ingress", SearchIngressApi)
	router.GET("/list/namespace", NamespaceApi)
	router.GET("/list/ingress", ListIngressApi)
	router.GET("/monitor/list",ListMonitor)
	router.GET("/montior/namespace",MonitorNamespace)
	return router
}

