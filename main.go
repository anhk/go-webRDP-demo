package main

import (
	"bytes"
	"fmt"
	"go-webRDP-demo/freerdp"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/gin-gonic/gin"
)

func rdpC() {
	//
	client := freerdp.NewClient("10.226.239.200",
		"administrator", "xThXxsP7mQ0Xufjux")

	if err := client.Connect(); err != nil {
		panic(err)
	}

	defer client.DisConnect()

	go func() {
		num := 0
		des := image.NewRGBA(image.Rect(0, 0, 1024, 768))
		_ = os.Mkdir("./tmp/", 0755)
		for {
			if msg, ok := client.Data(); !ok {
				break
			} else {
				img, err := png.Decode(bytes.NewReader(msg.Bitmap.Data))
				if err != nil {
					panic(err)
				}
				draw.Draw(des,
					image.Rect(msg.Bitmap.X, msg.Bitmap.Y, msg.Bitmap.X+msg.Bitmap.W, msg.Bitmap.Y+msg.Bitmap.H),
					img,
					image.Point{},
					draw.Over)
				if num++; num%100 == 0 {
					buf := new(bytes.Buffer)
					_ = png.Encode(buf, des)
					_ = os.WriteFile(fmt.Sprintf("./tmp/img-%v.png", num), buf.Bytes(), 0644)
				}
			}
		}
	}()
	if err := client.Start(); err != nil {
		panic(err)
	}
}

func main() {
	r := gin.Default()
	r.GET("/api/rdp", rdpProxy)
	r.Use(feMw("/"))
	_ = r.Run(":8081")
}
