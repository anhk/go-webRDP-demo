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
		for {
			var msg freerdp.Message
			if err := ws.ReadJSON(&msg); err != nil {
				fmt.Println("read from websocket fail:", err)
				client.DisConnect()
				break
			} else {
				if msg.Mouse != nil {
					client.ProcessMouseEvent(msg.Mouse)
				} else if msg.Keyboard != nil {
					client.ProcessKeyboardEvent(msg.Keyboard)
				}
			}
		}
	}()

	//cnt := 0
	//dest := image.NewRGBA(image.Rect(0, 0, 1024, 768))
	//_ = os.Mkdir("./tmp/", 0755)

	for {
		if msg, ok := client.Data(); !ok {
			fmt.Println(" !ok ")
			break
		} else if err := ws.WriteJSON(&msg); err != nil {
			fmt.Println(" err: ", err)
			break
		} else {
			// 打印图片
			//img, _ := png.Decode(bytes.NewReader(msg.Bitmap.Data))
			//draw.Draw(dest,
			//	image.Rect(msg.Bitmap.X, msg.Bitmap.Y, msg.Bitmap.X+msg.Bitmap.W, msg.Bitmap.Y+msg.Bitmap.H),
			//	img, image.Point{}, draw.Over)
			//if cnt++; cnt%100 == 0 {
			//	buf := new(bytes.Buffer)
			//	_ = png.Encode(buf, dest)
			//	_ = os.WriteFile(fmt.Sprintf("./tmp/img-%v.png", cnt), buf.Bytes(), 0644)
			//}

			// 打印信息
			//b64 := base64.StdEncoding.EncodeToString(msg.Bitmap.Data)
			//m := md5.Sum([]byte(b64))

			//fmt.Printf("send msg #%v, x:%v, y:%v, w:%v, h:%v, len:%v, md5:%x\n", cnt,
			//	msg.Bitmap.X, msg.Bitmap.Y, msg.Bitmap.W, msg.Bitmap.H, len(b64), m)
		}
	}
}
