package impex

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("impex", "GET", "models", restImpexListModels)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("impex", "GET", "export/:model", restImpexExportModel)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("impex", "POST", "import/:model", restImpexImportModel)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("impex", "POST", "import", restImpexImport)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("impex", "POST", "test/import", restImpexTestImport)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("impex", "POST", "test/mapping", restImpexTestCsvToMap)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API used to list available models for Impex system
func restImpexListModels(params *api.StructAPIHandlerParams) (interface{}, error) {
	var result []string

	for modelName, modelInstance := range models.GetDeclaredModels() {

		if _, present := impexModels[modelName]; present {
			continue
		}

		_, isObject := modelInstance.(models.InterfaceObject)
		_, isStorable := modelInstance.(models.InterfaceStorable)

		if isObject && isStorable {
			result = append(result, modelName)
		}
	}

	for modelName := range impexModels {
		result = append(result, modelName)
	}

	return result, nil
}

// WEB REST API used export specific model data from system
func restImpexExportModel(params *api.StructAPIHandlerParams) (interface{}, error) {

	modelName := params.RequestURLParams["model"]

	// if model, present := impexModels[modelName]; present {
	// 	model.Export()
	// }

	model, err := models.GetModel(modelName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	listable, isListable := model.(models.InterfaceListable)
	object, isObject := model.(models.InterfaceObject)

	if isListable && isObject {
		collection := listable.GetCollection()

		var attributes []string
		for _, attribute := range object.GetAttributesInfo() {
			attributes = append(attributes, attribute.Attribute)
			collection.ListAddExtraAttribute(attribute.Attribute)
		}

		// preparing csv writer
		csvWriter := csv.NewWriter(params.ResponseWriter)
		csvWriter.Comma = ','

		params.ResponseWriter.Header().Set("Content-type", "text/csv")
		params.ResponseWriter.Header().Set("Content-disposition", "attachment;filename=export_"+time.Now().Format(time.RFC3339)+".csv")

		// csvWriter.Write(attributes)
		// csvWriter.Flush()

		var records []map[string]interface{}

		list, _ := collection.List()
		for _, item := range list {
			records = append(records, item.Extra)

			/*record := make([]string, len(attributes))
			for idx, attribute := range attributes {
				record[idx] = utils.InterfaceToString(item.Extra[attribute])
			}
			csvWriter.Write(record)
			csvWriter.Flush()*/
		}

		MapToCSV(records, csvWriter)
	}

	return nil, nil
}

// WEB REST API used to test csv file before import
func restImpexTestImport(params *api.StructAPIHandlerParams) (interface{}, error) {

	params.ResponseWriter.Header().Set("Content-Type", "text/plain")

	filesProcessed := 0
	if params.Request != nil && params.Request.MultipartForm != nil && params.Request.MultipartForm.File != nil {
		for _, fileInfoArray := range params.Request.MultipartForm.File {
			for _, fileInfo := range fileInfoArray {

				attachedFile, err := fileInfo.Open()
				defer attachedFile.Close()
				if err != nil {
					return nil, env.ErrorDispatch(err)
				}

				// preparing csv reader
				csvReader := csv.NewReader(attachedFile)
				csvReader.Comma = ','

				err = ImportCSV(csvReader, params.ResponseWriter, true)
				if err != nil {
					return nil, env.ErrorDispatch(err)
				}
				filesProcessed++
			}
		}
	}

	return []byte(fmt.Sprintf("%d file(s) processed", filesProcessed)), nil
}

// WEB REST API used to process csv file script in impex format
func restImpexImport(params *api.StructAPIHandlerParams) (interface{}, error) {

	filesProcessed := 0
	if params.Request != nil && params.Request.MultipartForm != nil && params.Request.MultipartForm.File != nil {
		for _, fileInfoArray := range params.Request.MultipartForm.File {
			for _, fileInfo := range fileInfoArray {
				// if utils.IsAmongStr(fileInfo.Header.Get("Content-Type"), "application/csv", "text/csv") {
				attachedFile, err := fileInfo.Open()
				defer attachedFile.Close()
				if err != nil {
					return nil, env.ErrorDispatch(err)
				}

				// preparing csv reader
				csvReader := csv.NewReader(attachedFile)
				csvReader.Comma = ','

				err = ImportCSV(csvReader, params.ResponseWriter, false)
				if err != nil {
					return nil, env.ErrorDispatch(err)
				}

				filesProcessed++
			}
		}
	}

	return fmt.Sprintf("%d file(s) processed", filesProcessed), nil
}

// WEB REST API used import data to system
func restImpexImportModel(params *api.StructAPIHandlerParams) (interface{}, error) {

	modelName := params.RequestURLParams["model"]
	model, err := models.GetModel(modelName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	object, isObject := model.(models.InterfaceObject)
	_, isStorable := model.(models.InterfaceStorable)
	if !isObject || !isStorable {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e4a8bb15-7041-43ca-a29d-1c053409efef", modelName+" not implements InterfaceObject or InterfaceStorable interface")
	}

	attributes := make(map[string]models.StructAttributeInfo)
	for _, attribute := range object.GetAttributesInfo() {
		attributes[attribute.Attribute] = attribute
	}

	// start reading csv
	csvFile, _, err := params.Request.FormFile("file")
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ','

	// reading header
	csvColumns, err := csvReader.Read()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for _, csvColumn := range csvColumns {
		if _, ok := attributes[csvColumn]; !ok {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "de647436-9f79-45d8-80f9-139b52d1ba23", "there is no attribute "+csvColumn)
		}
	}

	for csvRecord, err := csvReader.Read(); err == nil; csvRecord, err = csvReader.Read() {
		model, _ = model.New()
		object, _ = model.(models.InterfaceObject)
		for idx, value := range csvRecord {
			object.Set(csvColumns[idx], value)
		}
		object.(models.InterfaceStorable).Save()
	}

	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "0a00d776-4158-43a2-8451-02f4344380d6", "not implemented")
}

// WEB REST API to test conversion from csv to json / map[string]interface{}
func restImpexTestCsvToMap(params *api.StructAPIHandlerParams) (interface{}, error) {
	var attachedFile interface {
		io.Reader
		io.Closer
	}
	var err error

	// looking for first one attached .csv file
	if params.Request != nil && params.Request.MultipartForm != nil && params.Request.MultipartForm.File != nil {

		for _, fileInfoArray := range params.Request.MultipartForm.File {
			for _, fileInfo := range fileInfoArray {
				if strings.Contains(fileInfo.Filename, ".csv") {

					attachedFile, err = fileInfo.Open()
					defer attachedFile.Close()
					if err != nil {
						return nil, env.ErrorDispatch(err)
					}

					break
				}
			}
		}

	}

	// file was not found
	if attachedFile == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6f5a582f-46f1-4b46-834d-e6b39da1ca68", "no csv file was attached")
	}

	// preparing csv reader
	var result []map[string]interface{}
	processor := func(item map[string]interface{}) bool {
		result = append(result, item)
		return true
	}

	reader := csv.NewReader(attachedFile)
	reader.Comma = ','

	// making csv processing
	exchangeDict := make(map[string]interface{})
	err = CSVToMap(reader, processor, exchangeDict)

	return result, err
}
