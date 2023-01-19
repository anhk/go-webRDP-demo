package main

import (
	"github.com/cihub/seelog"
	"go-webRDP-demo/freerdp"
)

func main() {
	defer seelog.Flush()

	client := freerdp.NewClient("10.226.239.200",
		"administrator", "xThXxsP7mQ0Xufjux")
	if err := client.Connect(); err != nil {
		panic(err)
	}

	defer client.DisConnect()

	if err := client.Start(); err != nil {
		panic(err)
	}
}
