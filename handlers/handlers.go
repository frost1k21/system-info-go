package handlers

import (
	. "example/generics/collectors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetSysInfoFromNames(c *gin.Context) {
	names := c.QueryArray("wsname")

	result := GetComputersInfo(names)
	c.JSON(http.StatusOK, result)
}

func PostSysInfoFromNames(c *gin.Context) {
	var wsNames []string
	err := c.BindJSON(&wsNames)
	if err != nil {
		log.Fatal(err.Error())
	}
	result := GetComputersInfo(wsNames)
	c.JSON(http.StatusOK, &result)
}
