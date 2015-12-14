package gallery

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	//	err = api.GetRestService().RegisterAPI("cms/gallery/image/:mediaName", api.ConstRESTOperationGet, APIGetGalleryImage)
	//	if err != nil {
	//		return env.ErrorDispatch(err)
	//	}

	err = api.GetRestService().RegisterAPI("cms/gallery/images", api.ConstRESTOperationGet, APIListGalleryImages)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("cms/gallery/images", api.ConstRESTOperationCreate, APIAddGalleryImages)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("cms/gallery/image/:mediaName", api.ConstRESTOperationDelete, APIRemoveGalleryImage)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	//	err = api.GetRestService().RegisterAPI("cms/gallery/path", api.ConstRESTOperationGet, APIGetGalleryPath)
	//	if err != nil {
	//		return env.ErrorDispatch(err)
	//	}

	return nil
}

// APIGetGalleryPath returns relative path to gallery library
//   - product id, media type must be specified in "productID" and "mediaType" arguments
func APIGetGalleryPath(context api.InterfaceApplicationContext) (interface{}, error) {

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	mediaPath, err := mediaStorage.GetMediaPath(ConstStorageModel, ConstStorageObject, ConstStorageType)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return mediaPath, nil
}

// APIListGalleryImages returns list of media files from gallery
func APIListGalleryImages(context api.InterfaceApplicationContext) (interface{}, error) {

	var result []interface{}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	mediaList, err := mediaStorage.ListMedia(ConstStorageModel, ConstStorageObject, ConstStorageType)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	path, err := mediaStorage.GetMediaPath(ConstStorageModel, ConstStorageObject, ConstStorageType)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	mediaBasePath := utils.InterfaceToString(env.ConfigGetValue("general.app.media_base_url"))
	path = mediaBasePath + "/" + path

	for _, mediaName := range mediaList {
		mediaNameParts := strings.SplitN(mediaName, ".", 2)
		creationDate := utils.InterfaceToTime(mediaNameParts[0][strings.LastIndex(mediaNameParts[0], "_")+1:])

		result = append(result, map[string]string{"name": mediaName, "url": path + mediaName, "created_at": creationDate.String()})
	}

	return result, nil
}

// APIAddGalleryImages uploads images to the gallery
//   - media file should be provided in "file" field with full name
func APIAddGalleryImages(context api.InterfaceApplicationContext) (interface{}, error) {
	var result []interface{}

	files := context.GetRequestFiles()
	if len(files) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "75a2ddaf-b63d-4eed-b16d-4b32778f5fc1", "media file was not specified")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	for fileName, fileReader := range files {
		fileContent, err := ioutil.ReadAll(fileReader)
		if err != nil {
			return result, env.ErrorDispatch(err)
		}

		if !strings.Contains(fileName, ".") {
			result = append(result, "Image: '"+fileName+"' should contain extension")
			continue
		}

		// Handle image name, adding unique values to name
		fileName = strings.TrimSpace(fileName)
		mediaNameParts := strings.SplitN(fileName, ".", 2)
		imageName := mediaNameParts[0] + utils.InterfaceToString(time.Now().Nanosecond()) + "_" + utils.InterfaceToString(time.Now().Unix()) + "." + mediaNameParts[1]

		// save to media storage operation
		err = mediaStorage.Save(ConstStorageModel, ConstStorageObject, ConstStorageType, imageName, fileContent)
		if err != nil {
			env.ErrorDispatch(err)
			result = append(result, "Image: '"+fileName+"' returned error on save")
			continue
		}

		result = append(result, fileName+": "+imageName)
	}

	return result, nil
}

// APIRemoveGalleryImage removes image from gallery
//   - media name must be specified in "mediaName" argument
func APIRemoveGalleryImage(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	imageName := context.GetRequestArgument("mediaName")
	if imageName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "63b37b08-3b21-48b7-9058-291bb7e635a1", "media name was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// remove media operation
	//---------------------
	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	err = mediaStorage.Remove(ConstStorageModel, ConstStorageObject, ConstStorageType, imageName)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APIGetGalleryImage returns image from gallery
//   - media name must be specified in "mediaName" argument
//   - on success case not a JSON data returns, but media file
func APIGetGalleryImage(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	imageName := context.GetRequestArgument("mediaName")
	if imageName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "124c8b9d-1a6b-491c-97ba-a03e8c828337", "media name was not specified")
	}

	// list media operation
	//---------------------
	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	imageFile, err := mediaStorage.Load(ConstStorageModel, ConstStorageObject, ConstStorageType, imageName)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}
	return imageFile, nil
}
