package seo

import (
	"errors"
	"os"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/utils"

	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/app/models/product"
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

	err = api.GetRestService().RegisterAPI("sitemap", "GET", "", restSitemapGenerate)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("sitemap", "GET", "sitemap.xml", restSitemap)
	if err != nil {
		return err
	}

	return nil
}

// WEB REST API function used to obtain url rewrites list
func restURLRewritesList(params *api.T_APIHandlerParams) (interface{}, error) {
	collection, err := db.GetCollection(COLLECTION_NAME_URL_REWRITES)
	if err != nil {
		return nil, err
	}

	collection.SetResultColumns("url", "type", "rewrite")
	records, err := collection.Load()

	return records, err
}

// WEB REST API function used to obtain rewrite for specified url
func restURLRewritesGet(params *api.T_APIHandlerParams) (interface{}, error) {
	collection, err := db.GetCollection(COLLECTION_NAME_URL_REWRITES)
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

	collection, err := db.GetCollection(COLLECTION_NAME_URL_REWRITES)
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
	collection, err := db.GetCollection(COLLECTION_NAME_URL_REWRITES)
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
	collection, err := db.GetCollection(COLLECTION_NAME_URL_REWRITES)
	if err != nil {
		return nil, err
	}

	err = collection.DeleteById(params.RequestURLParams["id"])
	if err != nil {
		return nil, err
	}

	return "ok", nil
}

// WEB REST API function used to get sitemap
func restSitemap(params *api.T_APIHandlerParams) (interface{}, error) {

	// if sitemap expied - generating new one
	info, err := os.Stat(SITEMAP_FILE_PATH)
	if err != nil || (time.Now().Unix()-info.ModTime().Unix() >= SITEMAP_EXPIRE_SEC) {
		return restSitemapGenerate(params)
	}

	// using generated otherwise
	sitemapFile, err := os.Open(SITEMAP_FILE_PATH)
	if err != nil {
		return nil, err
	}
	defer sitemapFile.Close()

	params.ResponseWriter.Header().Set("Content-Type", "text/xml")

	buffer := make([]byte, 1024)
	for {
		n, err := sitemapFile.Read(buffer)
		if err != nil || n == 0 {
			break
		}

		n, err = params.ResponseWriter.Write(buffer[0:n])
		if err != nil || n == 0 {
			break
		}
	}

	return nil, err
}

// WEB REST API function used to generate new sitemap
func restSitemapGenerate(params *api.T_APIHandlerParams) (interface{}, error) {

	params.ResponseWriter.Header().Set("Content-Type", "text/xml")

	// creating sitemap file
	sitemapFile, err := os.Create(SITEMAP_FILE_PATH)
	if err != nil {
		return nil, err
	}
	defer sitemapFile.Close()

	// writer to write into HTTP and file simultaneously
	newline := []byte("\n")
	writeLine := func(line []byte) {
		params.ResponseWriter.Write(line)
		params.ResponseWriter.Write(newline)

		sitemapFile.Write(line)
		sitemapFile.Write(newline)
	}

	// sitemap file preparations
	writeLine([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"))
	writeLine([]byte("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">"))

	baseURL := "http://dev.ottemo.io:8080/"
	pageType := ""

	// per database record iterator
	iteratorFunc := func(record map[string]interface{}) bool {
		pageUrl := ""
		if pageType == "" {
			pageUrl = baseURL + utils.InterfaceToString(record["url"])
		} else {
			pageUrl = baseURL + pageType + "/" + utils.InterfaceToString(record["_id"])
		}

		writeLine([]byte("  <url><loc>" + pageUrl + "</loc></url>"))

		return true
	}

	// Re-writed pages
	rewritesCollection, err := db.GetCollection(COLLECTION_NAME_URL_REWRITES)
	rewritesCollection.SetResultColumns("url")
	rewritesCollection.Iterate(iteratorFunc)

	rewritesCollection.SetResultColumns("rewrite")

	// Product pages
	pageType = "product"
	rewritesCollection.AddFilter("type", "=", pageType)

	productCollectionModel, _ := product.GetProductCollectionModel()
	dbProductCollection := productCollectionModel.GetDBCollection()
	dbProductCollection.SetResultColumns("_id")
	dbProductCollection.AddFilter("_id", "nin", rewritesCollection)
	dbProductCollection.Iterate(iteratorFunc)

	// Category pages
	pageType = "category"
	rewritesCollection.ClearFilters()
	rewritesCollection.AddFilter("type", "=", pageType)

	categoryCollectionModel, _ := category.GetCategoryCollectionModel()
	dbCategoryCollection := categoryCollectionModel.GetDBCollection()
	dbCategoryCollection.SetResultColumns("_id")
	dbCategoryCollection.AddFilter("_id", "nin", rewritesCollection)
	dbCategoryCollection.Iterate(iteratorFunc)

	// Cms pages
	pageType = "cms"
	rewritesCollection.ClearFilters()
	rewritesCollection.AddFilter("type", "=", pageType)

	cmsPageCollectionModel, _ := cms.GetCMSPageCollectionModel()
	dbCMSPageCollection := cmsPageCollectionModel.GetDBCollection()
	dbCMSPageCollection.SetResultColumns("_id")
	dbCMSPageCollection.AddFilter("_id", "nin", rewritesCollection)
	dbCMSPageCollection.Iterate(iteratorFunc)

	writeLine([]byte("</urlset>"))

	return nil, nil
}
