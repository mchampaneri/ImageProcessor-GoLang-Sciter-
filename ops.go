package main

import (
	"bytes"
	"encoding/base64"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

// Bright changes brightess of
// image in memory
func Bright(brightBy float64, thisDir string) string {
	img2 := imaging.AdjustBrightness(Images[Index], brightBy+NormalBirghtness)
	Images[Index] = img2
	mybuffer := new(bytes.Buffer)
	jpeg.Encode(mybuffer, img2, nil)
	return base64.StdEncoding.EncodeToString(mybuffer.Bytes())
}

// Sharpen changes sharpness of
// image in memory
func Sharpen(brightBy float64, thisDir string) string {
	img2 := imaging.Sharpen(Images[Index], brightBy+NormalBirghtness)
	Images[Index] = img2
	mybuffer := new(bytes.Buffer)
	jpeg.Encode(mybuffer, img2, nil)
	return base64.StdEncoding.EncodeToString(mybuffer.Bytes())
}
