package seo

import (
	"errors"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/utils"
)

func setupAPI() error {

	var err error = nil

	err = api.GetRestService().RegisterAPI("url_rewrite", "GET", "list", restURLRewritesList)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("url_rewrite", "GET", "get/:url", restURLRewritesGet)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("url_rewrite", "POST", "add", restURLRewritesAdd)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("url_rewrite", "PUT", "update/:id", restURLRewritesUpdate)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("url_rewrite", "DELETE", "delete/:id", restURLRewritesDelete)
	if err != nil {
		return err
	}

	return nil
}

// WEB REST API function used to obtain url rewrites list
func restURLRewritesList(params *api.T_APIHandlerParams) (interface{}, error) {
	collection, err := db.GetCollection(URL_REWRITES_COLLECTION_NAME)
	if err != nil {
		return nil, err
	}

	records, err := collection.Load()

	return records, err
}

// WEB REST API function used to obtain rewrite for specified url
func restURLRewritesGet(params *api.T_APIHandlerParams) (interface{}, error) {
	collection, err := db.GetCollection(URL_REWRITES_COLLECTION_NAME)
	if err != nil {
		return nil, err
	}

	collection.AddFilter("url", "=", params.RequestURLParams["url"])
	records, err := collection.Load()

	return records, err
}

// WEB REST API function used to update existing url rewrite
func restURLRewritesUpdate(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	postValues, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	collection, err := db.GetCollection(URL_REWRITES_COLLECTION_NAME)
	if err != nil {
		return nil, err
	}

	urlRewriteId := params.RequestURLParams["id"]
	record, err := collection.LoadById(urlRewriteId)
	if err != nil {
		return nil, err
	}

	// if rewrite 'url' was changed - checking new value for duplicates
	//-----------------------------------------------------------------
	if urlValue, present := postValues["url"]; present && urlValue != record["url"] {
		urlValue := utils.InterfaceToString(urlValue)

		collection.AddFilter("url", "=", urlValue)
		recordsNumber, err := collection.Count()
		if err != nil {
			return nil, err
		}
		if recordsNumber > 0 {
			return nil, errors.New("rewrite for url '" + urlValue + "' already exists")
		}

		record["url"] = urlValue
	}

	// updating other attributes
	//--------------------------
	attributes := []string{"type", "rewrite", "title", "meta_keywords", "meta_description"}
	for _, attribute := range attributes {
		if value, present := postValues[attribute]; present {
			record[attribute] = value
		}
	}

	// saving updates
	//---------------
	_, err = collection.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// WEB REST API function used to add url rewrite
func restURLRewritesAdd(params *api.T_APIHandlerParams) (interface{}, error) {

	// checking request params
	//------------------------
	postValues, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	if !utils.KeysInMapAndNotBlank(postValues, "url", "rewrite") {
		return nil, errors.New("'url' and 'rewrite' params should be specified")
	}

	valueUrl := utils.InterfaceToString(postValues["url"])
	valueRewrite := utils.InterfaceToString(postValues["rewrite"])

	// looking for duplicated 'url'
	//-----------------------------
	collection, err := db.GetCollection(URL_REWRITES_COLLECTION_NAME)
	if err != nil {
		return nil, err
	}

	collection.AddFilter("url", "=", valueUrl)
	recordsNumber, err := collection.Count()
	if err != nil {
		return nil, err
	}
	if recordsNumber > 0 {
		return nil, errors.New("rewrite for url '" + valueUrl + "' already exists")
	}

	// making new record and storing it
	//---------------------------------
	newRecord := map[string]interface{}{
		"url":              valueUrl,
		"type":             "",
		"rewrite":          valueRewrite,
		"title":            nil,
		"meta_keywords":    nil,
		"meta_description": nil,
	}

	attributes := []string{"type", "title", "meta_keywords", "meta_description"}
	for _, attribute := range attributes {
		if value, present := postValues[attribute]; present {
			newRecord[attribute] = value
		}
	}

	newId, err := collection.Save(newRecord)
	if err != nil {
		return nil, err
	}

	newRecord["_id"] = newId

	return newRecord, nil
}

// WEB REST API function used to delete url rewrite
func restURLRewritesDelete(params *api.T_APIHandlerParams) (interface{}, error) {
	collection, err := db.GetCollection(URL_REWRITES_COLLECTION_NAME)
	if err != nil {
		return nil, err
	}

	err = collection.DeleteById(params.RequestURLParams["id"])
	if err != nil {
		return nil, err
	}

	return "ok", nil
}
