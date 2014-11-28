package fsmedia

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"

	"io/ioutil"

	"image/jpeg"
	"image/png"

	"github.com/nfnt/resize"
	"github.com/ottemo/foundation/env"
)

// UpdateSizeNames loads predefined sizes from config value
func (it *FilesystemMediaStorage) UpdateSizeNames() {
	if configValue, ok := env.ConfigGetValue(ConstConfigPathMediaImageSizes).(string); ok && configValue != "" {
		for _, size := range strings.Split(configValue, ",") {
			var sizeName, sizeValue string
			sizeArray := strings.Split(size, ":")

			if len(sizeArray) > 1 {
				sizeName = strings.TrimSpace(sizeArray[1])
				sizeValue = strings.TrimSpace(sizeArray[0])
			} else {
				sizeValue = strings.TrimSpace(sizeArray[0])
			}

			if sizeValue != "" {
				if sizeName != "" {
					it.imageSizes[sizeName] = sizeValue
				}
				it.imageSizes[sizeValue] = sizeValue
			}

		}
	}
}

// GetSizeDimensions returns width and height for specified size or error if size is ot valid
func (it *FilesystemMediaStorage) GetSizeDimensions(size string) (uint, uint, error) {
	var height uint64

	if sizeValue, present := it.imageSizes[size]; present {
		size = sizeValue
	}

	value := strings.Split(size, "x")

	width, err := strconv.ParseUint(value[0], 10, 0)
	if err != nil {
		return 0, 0, errors.New("Invalid size")
	}

	if len(value) > 1 {
		height, _ = strconv.ParseUint(value[1], 10, 0)
	}

	return uint(width), uint(height), nil
}

// GetResizedMediaName returns media filename for a specified size, or error ir size is invalid
func (it *FilesystemMediaStorage) GetResizedMediaName(mediaName string, size string) string {

	var fileExtension string
	fileName := mediaName

	// checking file extension
	idx := strings.LastIndex(mediaName, ".")
	if idx != -1 {
		fileExtension = mediaName[idx:]
		fileName = mediaName[0:idx]
	}

	// if we have predefined size
	if _, present := it.imageSizes[size]; present {
		return fmt.Sprintf("%s_%s%s", fileName, size, fileExtension)
	}

	// otherwise
	width, height, _ := it.GetSizeDimensions(size)
	return fmt.Sprintf("%s_%dx%d", fileName, width, height)
}

// ResizeMediaImage re-sizes specified media to given size, returns error if not possible
func (it *FilesystemMediaStorage) ResizeMediaImage(model string, objID string, mediaName string, size string) error {
	path, _ := it.GetMediaPath(model, objID, "image")
	path = it.storageFolder + path

	sourceFileName := path + mediaName
	sourceImage, err := ioutil.ReadFile(sourceFileName)
	if err != nil {
		return err
	}

	width, height, err := it.GetSizeDimensions(size)
	if err != nil {
		return err
	}

	resizedFileName := path + it.GetResizedMediaName(mediaName, size)
	resizedFile, err := os.Create(resizedFileName)
	if err != nil {
		return err
	}

	decodedImage, imageFormat, err := image.Decode(bytes.NewReader(sourceImage))
	if err != nil {
		return err
	}

	originalSize := decodedImage.Bounds().Size()
	if width != 0 && height != 0 {
		if originalSize.X > originalSize.Y {
			height = 0
		} else {
			width = 0
		}
	}
	resizedImage := resize.Resize(width, height, decodedImage, resize.Bilinear)

	switch imageFormat {
	case "jpeg":
		err = jpeg.Encode(resizedFile, resizedImage, nil)
	case "png":
		err = png.Encode(resizedFile, resizedImage)
	default:
		return errors.New("unknown image format to encode")
	}

	if err != nil {
		return err
	}

	return nil
}

// GetResizedImage returns re-sized image contents or error if not possible
func (it *FilesystemMediaStorage) GetResizedImage(model string, objID string, mediaName string, size string) ([]byte, error) {
	path, _ := it.GetMediaPath(model, objID, "image")
	path = it.storageFolder + path

	resizedFileName := path + it.GetResizedMediaName(mediaName, size)
	resizedFileContents, err := ioutil.ReadFile(resizedFileName)
	if os.IsNotExist(err) {

		err = it.ResizeMediaImage(model, objID, mediaName, size)
		if err != nil {
			return nil, err
		}

		resizedFileContents, err = ioutil.ReadFile(resizedFileName)
		if err != nil {
			return nil, err
		}

	}

	return resizedFileContents, err
}
