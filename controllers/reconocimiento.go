package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/leandroveronezi/go-recognizer"
	"image"
	"image/png"
	"log"
	"os"
)

const dataDir = "utils"

func addFile(rec *recognizer.Recognizer, Path, Id string) {

	err := rec.AddImageToDataset(Path, Id)

	if err != nil {
		fmt.Println(err)
		return
	}

}

func Recog(ctx iris.Context) {

	rec := recognizer.Recognizer{}
	err := rec.Init(dataDir)

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	rec.Tolerance = 0.4
	rec.UseGray = true
	rec.UseCNN = false
	defer rec.Close()

	fileData, fileHeader, err := ctx.FormFile("file")
	imgSrc, _, err := image.Decode(fileData)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		log.Printf("failed decoding %s: %s", fileHeader.Filename, err)
		fmt.Println(err.Error())
	}
	file, err := os.Create(fileHeader.Filename)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		fmt.Println(err.Error())
		log.Printf("failed creating %s: %s", fileHeader.Filename, err)
		return
	}
	defer file.Close()
	png.Encode(file, imgSrc)

	pictureData, pictureHeader, err := ctx.FormFile("picture")
	picSrc, _, err := image.Decode(pictureData)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		log.Printf("failed decoding %s: %s", pictureHeader.Filename, err)
		fmt.Println(err.Error())
	}
	picture, err := os.Create(pictureHeader.Filename)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		fmt.Println(err.Error())
		log.Printf("failed creating %s: %s", pictureHeader.Filename, err)
		return
	}
	defer picture.Close()
	png.Encode(picture, picSrc)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	name := ctx.FormValue("name")
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	fmt.Println(picture.Name())
	fmt.Println(file.Name())

	addFile(&rec, picture.Name(), name)

	rec.SetSamples()

	faces, err := rec.ClassifyMultiples(file.Name())
	defer os.Remove(file.Name())
	defer os.Remove(picture.Name())

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	if len(faces) == 0 {
		ctx.StatusCode(iris.StatusOK)
		ctx.Header("Content-Type", "application/json")
		ctx.JSON(context.Map{"response": "No hubo coincidencia", "code": 01})
		return
	}

	//img, err := rec.DrawFaces(filepath.Join(fotosDir, "photoA.png"), faces)
	//
	//if err != nil {
	//	ctx.StatusCode(iris.StatusInternalServerError)
	//	fmt.Println(err.Error())
	//	return
	//}
	//
	//err = rec.SaveImage("controllers/faces.jpg", img)
	//if err != nil {
	//	ctx.StatusCode(iris.StatusInternalServerError)
	//	fmt.Println(err.Error())
	//	return
	//}

	ctx.StatusCode(iris.StatusOK)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(context.Map{"response": "Exitoso", "code": 00})
}
