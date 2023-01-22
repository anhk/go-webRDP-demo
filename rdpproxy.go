package main

import (
	"fmt"
	"go-webRDP-demo/freerdp"
	"net/http"

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

	client := freerdp.NewClient("192.168.1.129",
		"anhk", "my-password")

	if err := client.Connect(); err != nil {
		panic(err)
	}
	defer client.DisConnect()

	go func() {
		if err := client.Start(); err != nil {
			panic(err)
		}
	}()

	go func() {
		var msg freerdp.Message
		if err := ws.ReadJSON(&msg); err != nil {
			fmt.Println("read from websocket fail:", err)
			client.DisConnect()
		}
	}()

	for {
		if msg, ok := client.Data(); !ok {
			fmt.Println(" !ok ")
			break
		} else if err := ws.WriteJSON(&msg); err != nil {
			fmt.Println(" err: ", err)
			break
		}
	}
}
