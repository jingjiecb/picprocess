package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir := flag.String("d", ".", "Directory to search for image files")
	outputDir := flag.String("o", "./processed", "Output directory")

	if dir == outputDir {
		fmt.Println("Warning! Output directory should not be the same directory")
	}

	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	files, err := getImageFiles(*dir)
	if err != nil {
		fmt.Printf("unable to read pictures in %s, cause: %s\n", *dir, err)
		return
	}

	var errorFiles []string
	index := 1
	for _, file := range files {
		jpgImage, err := compressFileToJPEG(file, 300)
		if err != nil {
			fmt.Printf("fail to process all pictures in %s, cause: %s\n", *dir, err)
			continue
		}
		newJpgPath := getNewJpgPath(*outputDir, index)
		index++

		err = jpgImage.SaveToFile(newJpgPath)
		if err != nil {
			fmt.Printf("fail to save image to %s, cause: %s\n", newJpgPath, err)
			continue
		}

		fmt.Printf("Successfully processed: %s ==> %s", file, newJpgPath)
	}

	if len(errorFiles) > 0 {
		fmt.Printf("Some files cannot be processed, you can process them manually: %s\n", strings.Join(errorFiles, "; "))
	}
}

func getNewJpgPath(dir string, index int) string {
	return filepath.Join(dir, fmt.Sprintf("%d.jpg", index))
}

func getImageFiles(dir string) ([]string, error) {
	var imagePaths []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isImageFile(info.Name()) {
			imagePaths = append(imagePaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return imagePaths, nil
}

func isImageFile(fileName string) bool {
	extensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp"}
	ext := strings.ToLower(filepath.Ext(fileName))
	for _, e := range extensions {
		if ext == e {
			return true
		}
	}
	return false
}

func compressFileToJPEG(filePath string, maxKB int) (*JPGImage, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	if ok, jpgImage := tryGetJPGWithoutProcessing(file, maxKB); ok {
		return jpgImage, nil
	}

	file.Seek(0, 0)
	img, err := readImage(file)
	if err != nil {
		return nil, err
	}

	jpgImage, err := compressImageToJPEG(img, maxKB)
	if err != nil {
		return nil, err
	}
	return jpgImage, nil
}

func readImage(file *os.File) (image.Image, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func tryGetJPGWithoutProcessing(file *os.File, maxKB int) (bool, *JPGImage) {
	// if the image is already ok, skip
	jpgImage, err := NewJPGImageFromFile(file)
	if err == nil {
		if jpgImage.IsSizeLessOrEqualThan(maxKB) {
			return true, jpgImage
		}
	}

	return false, nil
}

func compressImageToJPEG(img image.Image, maxKB int) (*JPGImage, error) {
	quality := 100
	for quality > 1 {
		newJPGImage, err := NewJPGImage(img, quality)
		if err != nil {
			return nil, err
		}
		if newJPGImage.IsSizeLessOrEqualThan(maxKB) {
			return newJPGImage, nil
		}
		quality -= 5
	}

	return nil, fmt.Errorf("unable to compress image to %dKB", maxKB)
}
