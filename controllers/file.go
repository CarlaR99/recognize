package controllers

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/otiai10/gosseract/v2"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
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

func FileUpload4()(resp bool) {
	respData := true
	filename := "controllers/cedula.png"
	infile, err := os.Open(filename)

	if err != nil {
		log.Printf("failed opening %s: %s", filename, err)
		respData = false
		fmt.Println(err.Error())
	}
	defer infile.Close()
	newOffset, err := infile.Seek(0, 0)
	log.Printf(strconv.FormatInt(newOffset, 10))
	if err != nil {
		log.Printf("failed en la solucion %s: %s", filename, err)
		respData = false
		fmt.Println(err.Error())
	}
	imgSrc, _, err := image.Decode(infile)
	if err != nil {
		log.Printf("failed decoding %s: %s", filename, err)
		respData = false
		fmt.Println(err.Error())
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
		respData = false
		fmt.Println(err.Error())
	}
	defer newfile.Close()
	png.Encode(newfile, dstImage2)
	fmt.Println(respData)
	return respData
}

func FileUpload(ctx iris.Context) {

	if FileUpload4() {
		filename := "grayscale-cedula.png"
		infile, err := os.Open(filename)

		if err != nil {
			log.Printf("failed opening %s: %s", filename, err)
			fmt.Println(err.Error())
		}
		defer func(infile *os.File) {
			err := infile.Close()
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				fmt.Println(err.Error())
			}
		}(infile)

		// Create physical file
		tempfile, err := ioutil.TempFile("", "bancamiga"+"-")
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}
		defer func() {
			err := tempfile.Close()
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				fmt.Println(err.Error())
				return
			}
			err = os.Remove(tempfile.Name())
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				fmt.Println(err.Error())
				return
			}
		}()

		// Make uploaded physical
		if _, err = io.Copy(tempfile, infile); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}

		client := gosseract.NewClient()
		defer func(client *gosseract.Client) {
			err := client.Close()
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				fmt.Println(err.Error())
				return
			}
		}(client)

		err = client.SetImage(tempfile.Name())
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			fmt.Println(err.Error())
			return
		}
		client.Languages = []string{"spa"}
		//if langs := ctx.Request().FormValue("languages"); langs != "" {
		//	client.Languages = strings.Split(langs, ",")
		//}
		//if whitelist := ctx.Request().FormValue("whitelist"); whitelist != "" {
		//	client.SetWhitelist(whitelist)
		//}

		var out string
		switch ctx.Request().FormValue("format") {
		case "hocr":
			out, err = client.HOCRText()
		//render.EscapeHTML = false
		default:
			out, err = client.Text()
		}
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			fmt.Println(err.Error())
			return
		}
//		fmt.Println("GUENAS?")
//		fmt.Println(out)

		ctx.StatusCode(iris.StatusOK)
		ctx.Header("Content-Type", "application/json")
		ctx.JSON(context.Map{"response": strings.Trim(out, ctx.Request().FormValue("trim"))})
		//render.JSON(http.StatusOK, map[string]interface{}{
		//	"result":  strings.Trim(out, r.FormValue("trim")),
		//	"version": version,
		//})
	} else {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	// Get uploaded file
	//r.ParseMultipartForm(32 << 20)
	//// upload, h, err := r.FormFile("file")
	//upload, _, err := r.FormFile("file")
	//if err != nil {
	//	render.JSON(http.StatusBadRequest, err)
	//	return
	//}
	//defer upload.Close()

}
