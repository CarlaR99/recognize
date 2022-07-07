package controllers

import (
	"fmt"
	"html/template"
	"net/http"
)

// Index ...
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PASO??")
	fmt.Println(w)
	//var sherrifTmpl = template.New("index.html").Delims("{[{", "}]}")
	template.Must(template.ParseFiles("app/views/index.html")).ExecuteTemplate(w, "index.html", r)
	//marmoset.Render(w).HTML("index", map[string]interface{}{
	//	"AppName": "Recognize App",
	//})
}
