package main

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"io"
	"os"
)

type JPGImage struct {
	content *bytes.Buffer
}

func (jpg *JPGImage) SaveToFile(path string) error {
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = outFile.Write(jpg.content.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func NewJPGImage(image image.Image, quality int) (*JPGImage, error) {
	if quality <= 0 || quality > 100 {
		return nil, errors.New("quality must be between (0, 100]")
	}
	var content bytes.Buffer
	err := jpeg.Encode(&content, image, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}
	return &JPGImage{content: &content}, nil
}

func NewJPGImageFromFile(file *os.File) (*JPGImage, error) {
	isJPG, err := isJPEGUsingImage(file)
	if err != nil {
		return nil, err
	}

	if !isJPG {
		return nil, errors.New("image is not JPEG")
	}

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, file)
	return &JPGImage{content: &buffer}, nil
}

func (jpg *JPGImage) IsSizeLessOrEqualThan(sizeKB int) bool {
	return jpg.content.Len() <= sizeKB*1024
}

func isJPEGUsingImage(file *os.File) (bool, error) {
	_, format, err := image.Decode(file)
	if err != nil {
		return false, err
	}

	return format == "jpeg", nil
}
