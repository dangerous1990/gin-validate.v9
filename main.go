package main

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	route := gin.Default()
	initValidate()
	return route
}
func main() {
	setupRouter().Run(":8080")
}
