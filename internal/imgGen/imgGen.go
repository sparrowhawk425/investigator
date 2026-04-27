package main

import (
	"bufio"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"strings"

	"golang.org/x/image/draw"
)

// https://medium.com/@kushalchapagain123456/generating-ascii-art-of-image-in-golang-ef3badc71a74
func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Expected arg image file path")
	}
	fileBase := args[1]
	filePath := "assets/img/" + fileBase
	img, err := loadImage(filePath)
	if err != nil {
		log.Fatalf("error while accessing image: %v", err)
	}
	img = convGrayScale(resizeImage(img, 100))
	asciiLines := mapAscii(img)
	parts := strings.Split(fileBase, ".")
	outputFilePath := "assets/ascii/" + parts[0] + ".txt"
	err = writeAscii(asciiLines, outputFilePath)
	if err != nil {
		log.Fatalf("error writing output file %s: %v", outputFilePath, err)
	}
}

func loadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func resizeImage(img image.Image, width int) image.Image {
	bounds := img.Bounds()
	height := (bounds.Dy() * width) / bounds.Dx()
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(newImage, newImage.Bounds(), img, bounds, draw.Over, nil)
	return newImage
}

func convGrayScale(img image.Image) image.Image {
	bound := img.Bounds()
	grayImage := image.NewRGBA(bound)

	for i := bound.Min.X; i < bound.Max.X; i++ {
		for j := bound.Min.Y; j < bound.Max.Y; j++ {
			oldPixel := img.At(i, j)
			color := color.GrayModel.Convert(oldPixel)
			grayImage.Set(i, j, color)
		}
	}
	return grayImage
}

func mapAscii(img image.Image) []string {
	asciiChar := "$@B%#*+=,....."
	bound := img.Bounds()
	height, width := bound.Max.Y, bound.Max.X
	result := make([]string, height)

	for y := bound.Min.Y; y < height; y++ {
		var line strings.Builder
		for x := bound.Min.X; x < width; x++ {
			pixelValue := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			pixel := pixelValue.Y
			asciiIndex := int(pixel) * (len(asciiChar) - 1) / 255
			line.WriteString(string(asciiChar[asciiIndex]))
		}
		result[y] = line.String()
	}
	return result
}

func writeAscii(lines []string, outputFile string) error {

	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	w := bufio.NewWriter(outFile)
	for _, line := range lines {
		w.WriteString(line + "\n")
	}
	w.Flush()
	return nil
}
