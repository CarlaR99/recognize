package controllers

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract/v2"
	"github.com/otiai10/marmoset"
	"github.com/vitali-fedulov/images/v2"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	imgexp = regexp.MustCompile("^image")
	text   = ""
)

// FileUpload ...
//func FileUpload(w http.ResponseWriter, r *http.Request) {
//
//	render := marmoset.Render(w, true)
//
//	// Get uploaded file
//	r.ParseMultipartForm(32 << 20)
//	// upload, h, err := r.FormFile("file")
//	upload, _, err := r.FormFile("file")
//	if err != nil {
//		render.JSON(http.StatusBadRequest, err)
//		return
//	}
//	defer upload.Close()
//
//	// Create physical file
//	tempfile, err := ioutil.TempFile("", "ocrserver"+"-")
//	if err != nil {
//		render.JSON(http.StatusBadRequest, err)
//		return
//	}
//
//	defer func() {
//		tempfile.Close()
//		os.Remove(tempfile.Name())
//	}()
//	fmt.Println(tempfile.Name())
//
//	// Make uploaded physical
//	if _, err = io.Copy(tempfile, upload); err != nil {
//		render.JSON(http.StatusInternalServerError, err)
//		return
//	}
//
//	imageFile, _ := os.Open(tempfile.Name())
//	defer imageFile.Close()
//	imag, _, _ := image.Decode(imageFile)
//	//imgSrc, _, err := image.Decode(tempfile)
//
//	bounds := imag.Bounds()
//	y, h := bounds.Max.X, bounds.Max.Y
//	grayScale := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{y, h}})
//	for x := 0; x < y; x++ {
//		for y := 0; y < h; y++ {
//			imageColor := imag.At(x, y)
//			rr, gg, bb, _ := imageColor.RGBA()
//			r := math.Pow(float64(rr), 2.2)
//			g := math.Pow(float64(gg), 2.2)
//			b := math.Pow(float64(bb), 2.2)
//			m := math.Pow(0.2125*r+0.7154*g+0.0721*b, 1/2.2)
//			Y := uint16(m + 0.5)
//			grayColor := color.Gray{uint8(Y >> 8)}
//			grayScale.Set(x, y, grayColor)
//		}
//	}
//
//	//var img, _ = loadImage(tempfile.Name())
//	//var gray = rgbaToGray(img)
//
//	newFileName := "grayscale.png"
//	newfile, err := os.Create(newFileName)
//	if err != nil {
//		log.Printf("failed creating %s: %s", newfile, err)
//		panic(err.Error())
//	}
//	defer newfile.Close()
//	png.Encode(newfile, grayScale)
//
//	client := gosseract.NewClient()
//	defer client.Close()
//
//	client.SetImage(tempfile.Name())
//	client.Languages = []string{"eng"}
//	if langs := r.FormValue("languages"); langs != "" {
//		client.Languages = strings.Split(langs, ",")
//	}
//	if whitelist := r.FormValue("whitelist"); whitelist != "" {
//		client.SetWhitelist(whitelist)
//	}
//
//	var out string
//	switch r.FormValue("format") {
//	case "hocr":
//		out, err = client.HOCRText()
//		render.EscapeHTML = false
//	default:
//		out, err = client.Text()
//	}
//	if err != nil {
//		render.JSON(http.StatusBadRequest, err)
//		return
//	}
//
//	render.JSON(http.StatusOK, map[string]interface{}{
//		"result":  strings.Trim(out, r.FormValue("trim")),
//		"version": 1.0,
//	})
//}

func FileUpload3(w http.ResponseWriter, r *http.Request) {

	render := marmoset.Render(w, true)

	// Get uploaded file
	r.ParseMultipartForm(32 << 20)
	// upload, h, err := r.FormFile("file")
	upload, _, err := r.FormFile("file")
	if err != nil {
		render.JSON(http.StatusBadRequest, err)
		return
	}
	defer upload.Close()

	// Create physical file
	tempfile, err := ioutil.TempFile("", "bancamiga"+"-")
	if err != nil {
		render.JSON(http.StatusBadRequest, err)
		return
	}

	infile, err := os.Open(tempfile.Name())

	if err != nil {
		log.Printf("failed opening %s: %s", tempfile.Name(), err)
		panic(err.Error())
	}
	defer infile.Close()
	newOffset, err := infile.Seek(0, 0)
	log.Printf(strconv.FormatInt(newOffset, 10))
	if err != nil {
		log.Printf("failed en la solucion %s: %s", tempfile.Name(), err)
		panic(err.Error())
	}
	imgSrc, _, err := image.Decode(infile)
	if err != nil {
		log.Printf("failed decoding %s: %s", tempfile.Name(), err)
		panic(err.Error())
	}

	// Create a new grayscale image
	bounds := imgSrc.Bounds()
	yy, h := bounds.Max.X, bounds.Max.Y
	grayScale := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{yy, h}})
	for x := 0; x < yy; x++ {
		for y := 0; y < h; y++ {
			imageColor := imgSrc.At(x, y)
			rr, gg, bb, _ := imageColor.RGBA()
			r := math.Pow(float64(rr), 2.2)
			g := math.Pow(float64(gg), 2.2)
			b := math.Pow(float64(bb), 2.2)
			m := math.Pow(0.2125*r+0.7154*g+0.0721*b, 1/2.2)
			Y := uint16(m + 0.5)
			grayColor := color.Gray{uint8(Y >> 8)}
			grayScale.Set(x, y, grayColor)
		}
	}

	// Encode the grayscale image to the new file
	newFileName := "grayscale2.png"
	newfile, err := os.Create(newFileName)
	if err != nil {
		log.Printf("failed creating %s: %s", newfile, err)
		panic(err.Error())
	}
	defer newfile.Close()
	png.Encode(newfile, grayScale)

	defer func() {
		tempfile.Close()
		os.Remove(tempfile.Name())
	}()

	// Make uploaded physical
	if _, err = io.Copy(tempfile, newfile); err != nil {
		render.JSON(http.StatusInternalServerError, err)
		return
	}

	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(tempfile.Name())
	client.Languages = []string{"eng"}
	if langs := r.FormValue("languages"); langs != "" {
		client.Languages = strings.Split(langs, ",")
	}
	if whitelist := r.FormValue("whitelist"); whitelist != "" {
		client.SetWhitelist(whitelist)
	}

	var out string
	switch r.FormValue("format") {
	case "hocr":
		out, err = client.HOCRText()
		render.EscapeHTML = false
	default:
		out, err = client.Text()
	}
	if err != nil {
		render.JSON(http.StatusBadRequest, err)
		return
	}

	render.JSON(http.StatusOK, map[string]interface{}{
		"result":  strings.Trim(out, r.FormValue("trim")),
		"version": 1.0,
	})
}

func FileUpload2(car http.ResponseWriter, la *http.Request) {
	filename := "controllers/cedula.png"
	infile, err := os.Open(filename)
	render := marmoset.Render(car, true)

	if err != nil {
		log.Printf("failed opening %s: %s", filename, err)
		panic(err.Error())
	}
	defer infile.Close()
	newOffset, err := infile.Seek(0, 0)
	log.Printf(strconv.FormatInt(newOffset, 10))
	if err != nil {
		log.Printf("failed en la solucion %s: %s", filename, err)
		panic(err.Error())
	}
	imgSrc, _, err := image.Decode(infile)
	if err != nil {
		log.Printf("failed decoding %s: %s", filename, err)
		panic(err.Error())
	}

	// Create a new grayscale image
	bounds := imgSrc.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	grayScale := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			imageColor := imgSrc.At(x, y)
			rr, gg, bb, _ := imageColor.RGBA()
			r := math.Pow(float64(rr), 2.2)
			g := math.Pow(float64(gg), 2.2)
			b := math.Pow(float64(bb), 2.2)
			m := math.Pow(0.2125*r+0.7154*g+0.0721*b, 1/2.2)
			Y := uint16(m + 0.5)
			grayColor := color.Gray{uint8(Y >> 8)}
			grayScale.Set(x, y, grayColor)
		}
	}
	dstImage := imaging.AdjustContrast(grayScale, 40)
	dstImage2 := imaging.AdjustBrightness(dstImage, -20)

	// Encode the grayscale image to the new file
	newFileName := "grayscale-cedula.png"
	newfile, err := os.Create(newFileName)
	if err != nil {
		log.Printf("failed creating %s: %s", newfile, err)
		panic(err.Error())
	}
	defer newfile.Close()
	png.Encode(newfile, dstImage2)
	render.JSON(http.StatusOK, map[string]interface{}{
		"result":  newfile.Name(),
		"version": 1.0,
	})
}
func FileUpload(car http.ResponseWriter, la *http.Request) {
	render := marmoset.Render(car, true)
	// Open photos.
	imgA, err := images.Open("controllers/photoA.png")
	if err != nil {
		panic(err)
	}
	imgB, err := images.Open("controllers/photoB.png")
	if err != nil {
		panic(err)
	}

	// Calculate hashes and image sizes.
	hashA, imgSizeA := images.Hash(imgA)
	hashB, imgSizeB := images.Hash(imgB)
	// Image comparison.
	if images.Similar(hashA, hashB, imgSizeA, imgSizeB) {
		fmt.Println("Images are similar.")
		text = "Images are similar."
	} else {
		fmt.Println("Images are distinct.")
		text = "Images are distinct."
	}
	render.JSON(http.StatusOK, map[string]interface{}{
		"result":  text,
		"version": 1.0,
	})
}
