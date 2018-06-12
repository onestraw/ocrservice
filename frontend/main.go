package main

import (
	"flag"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

var (
	flagAddr        = flag.String("addr", ":10001", "listening address")
	flagBackendAddr = flag.String("backend.addr", "127.0.0.1:10002", "listening address")
)

func main() {
	flag.Parse()
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.LoadHTMLFiles("./app/views/index.html")
	router.Static("/assets", "./app/assets")
	router.GET("/", Index)
	router.GET("/status", Status)
	router.POST("/file", ReverseProxy())

	router.Run(*flagAddr)
}

func Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", map[string]interface{}{
		"AppName": "ocrservice",
	})
}

func Status(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "OK",
	})
}

// ReverseProxy submit the ocr task to backend service
func ReverseProxy() gin.HandlerFunc {
	target := *flagBackendAddr

	return func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = target
			req.URL.Path = "/ocrimage"
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
