package main

import (
	"fmt"
	"github.com/otiai10/ocrserver/controllers"
	"github.com/otiai10/ocrserver/utilities"
	"log"
	"net/http"
	"os"

	"github.com/otiai10/marmoset"
)

var logger *log.Logger

func main() {

	marmoset.LoadViews("./app/views")

	r := marmoset.NewRouter()
	// API
	r.GET("/status", utilities.Status)
	//r.POST("/base64", utilities.Base64)
	r.POST("/file", controllers.FileUpload)
	// Sample Page
	r.GET("/", controllers.Index)
	r.Static("/assets", "./app/assets")

	logger = log.New(os.Stdout, fmt.Sprintf("[%s] ", "ocrserver"), 0)
	r.Apply(&utilities.LogFilter{Logger: logger})

	//port := os.Getenv("PORT")
	port := "8080"
	if port == "" {
		logger.Fatalln("Required env `PORT` is not specified.")
	}
	logger.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Println(err)
	}
}
