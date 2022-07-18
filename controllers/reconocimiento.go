package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/leandroveronezi/go-recognizer"
	"path/filepath"
	"github.com/kataras/iris/v12/context"
)

const fotosDir = "controllers"
const dataDir = "models"

func addFile(rec *recognizer.Recognizer, Path, Id string) {

	err := rec.AddImageToDataset(Path, Id)

	if err != nil {
		fmt.Println(err)
		return
	}

}

func Recog(ctx iris.Context) {

	rec := recognizer.Recognizer{}
	err := rec.Init(fotosDir)

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		fmt.Println(err.Error())		
		return
	}

	rec.Tolerance = 0.4
	rec.UseGray = true
	rec.UseCNN = false
	defer rec.Close()

	addFile(&rec, filepath.Join(fotosDir, "photoB.png"), "Robert")

	rec.SetSamples()

	faces, err := rec.ClassifyMultiples(filepath.Join(fotosDir, "photoA.png"))

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
                fmt.Println(err.Error())		
		return
	}

	img, err := rec.DrawFaces(filepath.Join(fotosDir, "photoA.png"), faces)

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
                fmt.Println(err.Error())
		return
	}

	err = rec.SaveImage("controllers/faces.jpg", img)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
                fmt.Println(err.Error())		
		return 
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(context.Map{"response": "Exitoso"})
}
