package main

import (
	"github.com/gin-gonic/gin"
)

func WebServer() {
	r := gin.Default()
	r.GET("/api/rdp", rdpProxy)
	r.Use(feMw("/"))
	_ = r.Run(":8081")
}

func main() {
	WebServer()
}
