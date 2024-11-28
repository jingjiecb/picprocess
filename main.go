package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func main() {
	dir := flag.String("d", ".", "Directory to search for image files")
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

	var newFiles []string
	for _, file := range files {
		newFile, err := compressImageToJPEG(file, 300)
		if err != nil {
			fmt.Printf("fail to process all pictures in %s, cause: %s\n", *dir, err)
		}
		newFiles = append(newFiles, newFile)
	}

	err = renameFiles(newFiles)
	if err != nil {
		fmt.Printf("fail to rename files, cause: %s\n", err)
		return
	}

	fmt.Printf("%d files processed.\n", len(newFiles))
}

func renameFiles(filePaths []string) error {
	for i, filePath := range filePaths {
		ext := filepath.Ext(filePath)
		newFileName := fmt.Sprintf("%d%s", i+1, ext)

		dir := filepath.Dir(filePath)

		newFilePath := filepath.Join(dir, newFileName)

		err := os.Rename(filePath, newFilePath)
		if err != nil {
			return fmt.Errorf("failed to rename file %s to %s: %v", filePath, newFilePath, err)
		}

		fmt.Printf("Renamed: %s -> %s\n", filePath, newFilePath)
	}

	return nil
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

func compressImageToJPEG(inputPath string, maxKB int) (string, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return inputPath, fmt.Errorf("failed to open file %s: %v", inputPath, err)
	}

	var img image.Image
	originalSupportedExt := []string{".jpg", ".jpeg", ".png", ".gif"}
	ext := filepath.Ext(inputPath)
	if slices.Contains(originalSupportedExt, ext) {
		img, _, err = image.Decode(file)
	} else {
		err = fmt.Errorf("unsupported file extension: %s", ext)
	}
	if err != nil {
		return inputPath, fmt.Errorf("failed to decode image %s: %v", inputPath, err)
	}

	err = file.Close()
	if err != nil {
		fmt.Printf("failed to close file %s: %v", inputPath, err)
		return inputPath, err
	}

	var buffer bytes.Buffer
	quality := 100
	for quality >= 10 {
		buffer.Reset()

		err = jpeg.Encode(&buffer, img, &jpeg.Options{Quality: quality})
		if err != nil {
			return "", fmt.Errorf("failed to encode image: %v", err)
		}

		if buffer.Len() <= maxKB*1024 {
			break
		}
		quality -= 5
	}

	if buffer.Len() > maxKB*1024 {
		return inputPath, fmt.Errorf("unable to compress image %s to %dKB", inputPath, maxKB)
	}

	fileName := filepath.Base(inputPath)
	filePosition := filepath.Dir(inputPath)
	outputPath := filepath.Join(filePosition, strings.TrimSuffix(fileName, filepath.Ext(fileName))+".jpg")

	_ = os.Remove(outputPath)
	_ = os.Remove(inputPath)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return outputPath, fmt.Errorf("failed to create output file %s: %v", inputPath, err)
	}
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {
			fmt.Printf("failed to close output file %s: %v", inputPath, err)
		}
	}(outFile)

	_, err = outFile.Write(buffer.Bytes())
	if err != nil {
		return inputPath, fmt.Errorf("failed to write output file %s: %v", inputPath, err)
	}

	return outputPath, nil
}
