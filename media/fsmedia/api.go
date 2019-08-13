package fsmedia

import (
	"github.com/ottemo/commerce/api"
	"github.com/ottemo/commerce/media"
	"github.com/ottemo/commerce/utils"
	"strings"
)

// configures package related API endpoint routines
func setupAPI() error {

	// api.GetRestService().GET("media/resizeAll", api.IsAdminHandler(APIResizeAll))
	api.GetRestService().GET("media/*path", APIGetMedia)

	return nil
}

func APIGetMedia(context api.InterfaceApplicationContext) (interface{}, error)  {
	path := context.GetRequestArgument("path")

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return nil, err
	}

	result, err := mediaStorage.GetMediaByPath(path)
	if err != nil {
		return nil, err
	}

	// ref to https://www.sitepoint.com/mime-types-complete-list/
	if strings.HasSuffix(path, ".jpg") {
		context.SetResponseContentType("image/jpeg")
	} else if strings.HasSuffix(path, ".png") {
		context.SetResponseContentType("image/png")
	} else if strings.HasSuffix(path, ".avi") {
		context.SetResponseContentType("image/avi")
	} else if strings.HasSuffix(path, ".txt") {
		context.SetResponseContentType("text/plain")
	} else {
		context.SetResponseContentType("application/octet-stream")
	}

	return result, nil
}

// APIGetMediaInfo will resize all images if the params of the request contain 'resizeAll' with a value of true
func APIResizeAll(context api.InterfaceApplicationContext) (interface{}, error) {
	// TODO: add example api call or add this to Apiary - jwv

	requestParams := context.GetRequestArguments()
	resizeAll := utils.GetFirstMapValue(requestParams, "resizeAll", "resizeImages", "resizeAllImages")

	if resizeAll != nil && utils.InterfaceToBool(resizeAll) {
		mediaStorage, err := media.GetMediaStorage()
		if err != nil {
			return nil, err
		}

		err = mediaStorage.ResizeAllMediaImages()
		if err != nil {
			return nil, err
		}
	}

	return "ok", nil
}
