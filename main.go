package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/disintegration/imaging"
	"github.com/fatih/color"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

const (
	NormalBirghtness = -50.0
)

var Index int            // Stores current index of image
var Images []image.Image // Images stores base64 string of images
var Files []os.FileInfo

func main() {
	// make rect for window
	rect := sciter.NewRect(0, 0, 800, 600)

	// create a window using upper rect
	win, _ := window.New(sciter.SW_MAIN|sciter.SW_ENABLE_DEBUG, rect)

	win.SetTitle("ImageViewar+-")

	// Scanning and loading base64 of images
	findAndLoadImageFromCurrentDirectory()

	// registering methods
	win.DefineFunction("loadFirstImage", LoadFirstImage)
	win.DefineFunction("loadNextImage", LoadNextImage)
	win.DefineFunction("loadPreviousImage", LoadPreviousImage)
	win.DefineFunction("brightCurrentImage", brightCurrentImage)
	win.DefineFunction("closeWindow", closeApplication)

	// Getting data from archive
	win.SetResourceArchive(resources)
	win.LoadFile("this://app/htdocs/image-viewer.htm")

	// Running app
	win.Show()
	win.Run()
	win.CloseArchive()
}

func closeApplication(vals ...*sciter.Value) *sciter.Value {
	syscall.Exit(0)
	return nil
}

// findAndLoadImageFromCurrentDirectory scans directory
// in which exec is for jpg / png files. Reads those files
// and store base64 string of those file in Images([]string)
func findAndLoadImageFromCurrentDirectory() {

	var waitGroup sync.WaitGroup
	// Getting working directory
	thisDir, dirErr := os.Getwd()
	if dirErr != nil {
		fmt.Println("Failed to get current directory")
		return
	}
	files, readDirErr := ioutil.ReadDir(thisDir)
	if readDirErr != nil {
		fmt.Println("failed to read current directory")
		return
	}

	if len(files) > 0 {
		img := getImage(files[0], thisDir)
		if img != nil {
			color.Yellow("First image has been loaded")
			Images = append(Images, img)
			Files = append(Files, files[0])
		}
	}

	// Loading files excpet first via goroutine
	// so we don't have to wait for every image
	// to be loaded to show up first image
	waitGroup.Add(1)
	go func() {
		for i, file := range files {
			if i == 0 {
				continue
			}
			img := getImage(file, thisDir)
			if img != nil {
				Images = append(Images, img)
				Files = append(Files, file)
			}
		}
		waitGroup.Done()
	}()
	waitGroup.Wait()
}

// LoadFirstImage return first
// image from Image array
// to sciter
func LoadFirstImage(vals ...*sciter.Value) *sciter.Value {
	if len(Images) > 0 {
		Index = 0
		buf := new(bytes.Buffer)
		png.Encode(buf, Images[0])
		return sciter.NewValue(base64.StdEncoding.EncodeToString(buf.Bytes()))
	}
	return sciter.NewValue(string("-"))
}

// LoadNextImage return image from
// next index to current index
func LoadNextImage(vals ...*sciter.Value) *sciter.Value {
	if Index < len(Images)-1 {
		Index++
		buf := new(bytes.Buffer)
		png.Encode(buf, Images[Index])
		return sciter.NewValue(base64.StdEncoding.EncodeToString(buf.Bytes()))
	}
	return LoadFirstImage()
}

// LoadPreviousImage return image from
// previous index to current index
func LoadPreviousImage(vals ...*sciter.Value) *sciter.Value {
	if Index > 0 {
		Index--
		buf := new(bytes.Buffer)
		png.Encode(buf, Images[Index])
		return sciter.NewValue(base64.StdEncoding.EncodeToString(buf.Bytes()))
	}
	return LoadFirstImage()
}

func brightCurrentImage(vals ...*sciter.Value) *sciter.Value {
	cwd, _ := os.Getwd()
	fmt.Println("your brightness perameter is ", vals[0].Float())
	imageString := Bright(vals[0].Float(), Files[Index], cwd)
	thisString := sciter.NewValue(imageString)
	return thisString
}

// getImage returns base64 string
// of file provided as input
func getImage(file os.FileInfo, thisDir string) image.Image {

	// Just supporting jpg and png file to be loaded
	// others are on the way .. .
	if strings.Contains(file.Name(), ".jpg") || strings.Contains(file.Name(), ".png") {
		imageFile, imageFileErr := os.Open(filepath.Join(thisDir, file.Name()))
		if imageFileErr != nil {
			fmt.Println("failed to load image file")
			return nil
		}

		// Reading image file in buffer
		fReader := bufio.NewReader(imageFile)

		Image, _, imageReadingErr := image.Decode(fReader)
		if imageReadingErr != nil {
			fmt.Println("faild to read image in image.Image", imageReadingErr.Error())
			return nil
		}
		// Convert file to base64
		// imgStrging := base64.StdEncoding.EncodeToString(buf)
		return Image
	}
	return nil
}

func Bright(brightBy float64, file os.FileInfo, thisDir string) string {

	img2 := imaging.AdjustBrightness(Images[Index], brightBy+NormalBirghtness)
	mybuffer := new(bytes.Buffer)
	jpeg.Encode(mybuffer, img2, nil)
	return base64.StdEncoding.EncodeToString(mybuffer.Bytes())

}

// BlurImage
// func Blur(file os.FileInfo, thisDir string) string {
// 	fmt.Println("blurring image")
// 	imageFile, imageFileErr := os.Open(filepath.Join(thisDir, file.Name()))
// 	if imageFileErr != nil {
// 		fmt.Println("failed to load image file")
// 		return ""
// 	}
// 	srcImage, _, err := image.Decode(imageFile)
// 	if err != nil {
// 		fmt.Println("failed to load decode image")
// 		return ""
// 	}
// 	dstImage := imaging.Blur(srcImage, 0.5)
// 	tempDir := os.TempDir()
// 	name, _ := uuid.NewV4()
// 	tempFile, errTemp := os.OpenFile(path.Join(tempDir, name.String()),
// 		os.O_CREATE|os.O_RDWR, os.ModeTemporary)

// 	if errTemp != nil {
// 		fmt.Println("failed to create temp file to store image ")
// 		return ""
// 	}
// 	encodingFialed := png.Encode(tempFile, dstImage)
// 	if encodingFialed != nil {
// 		fmt.Println("failed to encode file to return", encodingFialed.Error())
// 		return ""
// 	}

// 	state, statError := tempFile.Stat()
// 	if statError != nil {
// 		fmt.Println("failed to get error state ", statError.Error())
// 		return ""
// 	}

// 	size := state.Size()
// 	buf := make([]byte, size)

// 	// Reading image file in buffer
// 	fReader := bufio.NewReader(imageFile)
// 	fReader.Read(buf)
// 	imgStrging := base64.StdEncoding.EncodeToString(buf)
// 	fmt.Println(imgStrging)
// 	return imgStrging
// }
