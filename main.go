package main

import (
	. "example/generics/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()

	router.GET("/sysinfos", GetSysInfoFromNames)
	router.POST("/sysinfos", PostSysInfoFromNames)

	log.Fatal(router.Run("localhost:8080"))
}
