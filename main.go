package main

import (
	"github.com/gin-gonic/gin"
)

func WebServer() {
	//gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.GET("/api/rdp", rdpProxy)
	r.Use(feMw("/"))
	_ = r.Run(":8082")
}

func main() {
	WebServer()
}
