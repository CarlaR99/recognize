package controllers

import (
	"fmt"
	"github.com/leandroveronezi/go-recognizer"
	"net/http"
	"path/filepath"
)

const fotosDir = "fotos"
const dataDir = "models"

func addFile(rec *recognizer.Recognizer, Path, Id string) {

	err := rec.AddImageToDataset(Path, Id)

	if err != nil {
		fmt.Println(err)
		return
	}

}

func Recog(w http.ResponseWriter, r *http.Request) {

	rec := recognizer.Recognizer{}
	err := rec.Init(dataDir)

	if err != nil {
		fmt.Println(err)
		return
	}

	rec.Tolerance = 0.4
	rec.UseGray = true
	rec.UseCNN = false
	defer rec.Close()

	addFile(&rec, filepath.Join(fotosDir, "controllers/photoB.png"), "Robert")
	//addFile(&rec, filepath.Join(fotosDir, "bernadette.jpg"), "Bernadette")
	//addFile(&rec, filepath.Join(fotosDir, "howard.jpg"), "Howard")
	//addFile(&rec, filepath.Join(fotosDir, "penny.jpg"), "Penny")
	//addFile(&rec, filepath.Join(fotosDir, "raj.jpg"), "Raj")
	//addFile(&rec, filepath.Join(fotosDir, "sheldon.jpg"), "Sheldon")
	//addFile(&rec, filepath.Join(fotosDir, "leonard.jpg"), "Leonard")

	rec.SetSamples()

	faces, err := rec.ClassifyMultiples(filepath.Join(fotosDir, "controllers/photoA.png"))

	if err != nil {
		fmt.Println(err)
		return
	}

	img, err := rec.DrawFaces(filepath.Join(fotosDir, "controllers/photoA.png"), faces)

	if err != nil {
		fmt.Println(err)
		return
	}

	rec.SaveImage("controllers/faces.jpg", img)

}
