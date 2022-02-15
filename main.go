package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

// Symbol - created to organise unique symbols for project.
type Symbol string

const (
	X Symbol = "X"
	O Symbol = "O"
	H Symbol = "H"
)

const totalImages int = 20

func main() {
	var wg sync.WaitGroup

	for y := 1; y <= totalImages; y++ {
		wg.Add(1)
		go generateImage(y, &wg)
	}

	wg.Wait()
	fmt.Println("Successfully created " + strconv.Itoa(totalImages) + " images.")
}

func generateImage(currentNo int, wg *sync.WaitGroup) {
	defer wg.Done()

	white := color.RGBA{R: 255, G: 255, B: 255, A: 1}

	canvas := image.NewRGBA(image.Rect(0, 0, 500, 500))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: white}, image.Point{}, draw.Src)

	for i := 1; i <= 9; i++ {

		// Generate a random seed and times it by 4 to create some extra randomness.
		rand.Seed(time.Now().UnixNano() * 4)
		random := rand.Intn(1000 - 1)

		if random >= 900 {
			index := strconv.Itoa(i)
			x := getImage(Symbol(X), index)
			draw.Draw(canvas, canvas.Bounds(), x, image.Point{}, draw.Over)
		} else if random < 900 && random >= 600 {
			index := strconv.Itoa(i)
			o := getImage(Symbol(O), index)
			draw.Draw(canvas, canvas.Bounds(), o, image.Point{}, draw.Over)
		} else {
			index := strconv.Itoa(i)
			h := getImage(Symbol(H), index)
			draw.Draw(canvas, canvas.Bounds(), h, image.Point{}, draw.Over)
		}
	}

	var imageNo = strconv.Itoa(currentNo)
	toImg, err := os.Create("./export/output-" + imageNo + ".jpg")
	if err != nil {
		panic(err)
	}

	defer toImg.Close()

	err = jpeg.Encode(toImg, canvas, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Saved output-" + imageNo + ".jpg")

	// Refresh Canvas
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{white}, image.Point{0, 0}, draw.Src)
}

func getImage(symbol Symbol, num string) image.Image {
	file, err := os.Open("./images/" + string(symbol) + num + ".png")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	decodedFile, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	return decodedFile
}
