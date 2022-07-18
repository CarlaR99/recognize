package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
)

// Index ...
func Index(ctx iris.Context) {
	fmt.Println("PASO??")
//	fmt.Println(w)
	err := ctx.View("index.html")
	if err != nil {
		return
	}

//	template.Must(template.ParseFiles("app/views/index.html")).ExecuteTemplate(w, "index.html", r)
	//marmoset.Render(w).HTML("index", map[string]interface{}{
	//	"AppName": "Recognize App",
	//})
}
