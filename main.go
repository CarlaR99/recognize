package main

import (
	"github.com/kataras/iris/v12"
	"github.com/otiai10/ocrserver/controllers"
)

func main() {

	app := iris.New()
	tmpl := iris.HTML("./app/views", ".html")
	app.RegisterView(tmpl)

	app.Get("/", controllers.Index)
	app.Post("/file", controllers.FileUpload)
	app.Post("/recognize", controllers.Recog)
	app.HandleDir("/assets", "./app/assets")
	
	err := app.Run(iris.Addr(":5000"))
	if err != nil {
		return
	}

}
