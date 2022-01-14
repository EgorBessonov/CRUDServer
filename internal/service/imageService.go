package service

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"os"
)

// UploadImage method save image in local folder
func (s Service) UploadImage(image *multipart.FileHeader) error {
	src, err := image.Open()
	if err != nil {
		return fmt.Errorf("image service: can't upload image - %e", err)
	}
	defer func() {
		err := src.Close()
		if err != nil {
			log.Error("error while closing multipart file instance.")
		}
	}()
	dst, err := os.Create("images/" + image.Filename)
	if err != nil {
		return fmt.Errorf("image service: can't upload image - %e", err)
	}
	defer func() {
		err := dst.Close()
		if err != nil {
			log.Error("error while closing os file instance.")
		}
	}()
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("image service: can't upload image - %e", err)
	}
	return nil
}
