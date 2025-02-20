package utils

import (
	"bytes"
	"image/jpeg"
	"image/png"
)

func ResizeImage(imgData []byte, maximumSize int) ([]byte, error) {
	return imgData, nil
}

func ConvertPNGToJPG(imgData []byte, quality int) ([]byte, error) {
	// Create a bytes.Reader from imgData
	imgReader := bytes.NewReader(imgData)

	// Decode the PNG image
	img, err := png.Decode(imgReader)
	if err != nil {
		return nil, err
	}

	// Create a bytes.Buffer to write the JPEG image to
	jpgBuffer := new(bytes.Buffer)

	// Encode the image to JPEG
	err = jpeg.Encode(jpgBuffer, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}

	// Return the JPEG image as a byte slice
	return jpgBuffer.Bytes(), nil
}
