package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
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

const totalImages int = 100

func main() {
	fmt.Println("Generating " + strconv.Itoa(totalImages) + " images...")
	fmt.Println("Running...")
	var wg sync.WaitGroup

	for y := 1; y <= totalImages; y++ {
		wg.Add(1)
		go generateImage(y, &wg)
	}

	wg.Wait()
	fmt.Println("Successfully created " + strconv.Itoa(totalImages) + " images.")

	files, err := ioutil.ReadDir("./export")
	if err != nil {
		panic(err)
	}

	// Map all files to hashmap with amount of times they're found.
	filesMap := make(map[string]int)
	shaFilenameMap := make(map[string]string)
	for _, f := range files {
		sha256Struct := sha256.New()

		openedFile, err := os.Open("./export/" + f.Name())
		if err != nil {
			panic(err)
		}

		_, err = io.Copy(sha256Struct, openedFile)
		if err != nil {
			panic(err)
		}

		decodedSha256 := base64.URLEncoding.EncodeToString(sha256Struct.Sum(nil))

		filesMap[decodedSha256]++
		shaFilenameMap[decodedSha256] = f.Name()
	}

	// Compare checksums to find duplicates.
	deletedCount := 0
	var toDelete []string
	for key, element := range filesMap {
		if element >= 2 {
			toDelete = append(toDelete, shaFilenameMap[key])

			// Counting the amount of images removed.
			deletedCount++
		}
	}

	// Remove duplicates here.
	for _, element := range toDelete {
		err := os.Remove("./export/" + element)
		if err != nil {
			panic(err)
		}
	}

	files, err = ioutil.ReadDir("./export")
	if err != nil {
		panic(err)
	}

	// Update file names to not have gaps.
	for i, f := range files {
		newFileName := "./export/output" + strconv.Itoa(i) + ".jpg"
		err := os.Rename("./export/"+f.Name(), newFileName)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Successfully deleted " + strconv.Itoa(deletedCount) + " duplicates.")
	fmt.Println("Successfully reorganised output file names.")
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

	if currentNo%100 == 0 {
		fmt.Println("output-" + strconv.Itoa(currentNo) + ".jpg ++")
	}

	// Refresh Canvas
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: white}, image.Point{}, draw.Src)
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
