package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

////go:embed static/*
//var dist embed.FS

const fsBase = "static"

func feMw(urlPrefix string) gin.HandlerFunc {
	const indexHtml = "index.html"

	return func(c *gin.Context) {
		urlPath := strings.TrimSpace(c.Request.URL.Path)
		if urlPath == urlPrefix {
			urlPath = path.Join(urlPrefix, indexHtml)
		}
		urlPath = filepath.Join(fsBase, urlPath)
		f, err := os.Open(urlPath)
		//f, err := dist.Open(urlPath)
		if err != nil {
			return
		}
		fi, err := f.Stat()
		if strings.HasSuffix(urlPath, ".html") {
			c.Header("Cache-Control", "no-cache")
			c.Header("Content-Type", "text/html; charset=utf-8")
		}

		if strings.HasSuffix(urlPath, ".js") {
			c.Header("Content-Type", "text/javascript; charset=utf-8")
		}
		if strings.HasSuffix(urlPath, ".css") {
			c.Header("Content-Type", "text/css; charset=utf-8")
		}

		if err != nil || !fi.IsDir() {
			//bs, err := dist.ReadFile(urlPath)
			bs, err := os.ReadFile(urlPath)
			if err != nil {
				return
			}
			c.Status(200)
			_, _ = c.Writer.Write(bs)
			c.Abort()
		}
	}
}
