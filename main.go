package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"currency-gin/controller"
	"os"
	"io"
)

func main() {
	gin.DisableConsoleColor()

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/convert/", new(controller.ConverterController).Handle)
	r.GET("/api/convert/:fromUnit", new(controller.ConverterController).Handle)
	r.GET("/api/convert/:fromUnit/:toUnit", new(controller.ConverterController).Handle)

	r.GET("/api/youtube", new(controller.YoutubeController).Handle)
	r.POST("/api/youtube", new(controller.YoutubeController).Handle)
	r.GET("/api/youtube/:method", new(controller.YoutubeController).Handle)
	r.POST("/api/youtube/:method", new(controller.YoutubeController).Handle)
	r.Static("/youtube", "./static/public/youtube")

	//r.Run() // listen and serve on 0.0.0.0:8080

	s := &http.Server{
		Addr:			":8080",
		Handler:		r,
		ReadTimeout:	10 * time.Second,
		WriteTimeout:	10 * time.Second,
		MaxHeaderBytes:	1 << 20,
	}
	s.ListenAndServe()

}
