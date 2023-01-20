package main

import (
	"go-webRDP-demo/freerdp"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func rdpProxy(ctx *gin.Context) {
	proto := ctx.Request.Header.Get("Sec-Websocket-Protocol")

	ws, err := upgrade.Upgrade(ctx.Writer, ctx.Request, http.Header{
		"Sec-Websocket-Protocol": {proto},
	})
	if err != nil {
		panic(err)
	}

	for {
		msg := freerdp.Message{}
		_ = ws.WriteJSON(&msg)
		time.Sleep(time.Second)
	}
}
