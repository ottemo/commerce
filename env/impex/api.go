package impex

import (
	"encoding/csv"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func setupAPI() error {

	var err error = nil

	err = api.GetRestService().RegisterAPI("impex", "GET", "export/:model", restImpexExport)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("impex", "POST", "import/:model", restImpexImport)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API used export specific model data from system
func restImpexExport(params *api.T_APIHandlerParams) (interface{}, error) {

	model, err := models.GetModel(params.RequestURLParams["model"])
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	listable, isListable := model.(models.I_Listable)
	object, isObject := model.(models.I_Object)

	if isListable && isObject {
		collection := listable.GetCollection()

		attributes := make([]string, 0)
		for _, attribute := range object.GetAttributesInfo() {
			attributes = append(attributes, attribute.Attribute)
			collection.ListAddExtraAttribute(attribute.Attribute)
		}

		// preparing csv writer
		csvWriter := csv.NewWriter(params.ResponseWriter)

		params.ResponseWriter.Header().Set("Content-type", "text/csv")
		params.ResponseWriter.Header().Set("Content-disposition", "attachment;filename=export_"+time.Now().Format(time.RFC3339)+".csv")

		csvWriter.Write(attributes)
		csvWriter.Flush()

		list, _ := collection.List()
		for _, item := range list {
			record := make([]string, len(attributes))
			for idx, attribute := range attributes {
				record[idx] = utils.InterfaceToString(item.Extra[attribute])
			}
			csvWriter.Write(record)
			csvWriter.Flush()
		}

	}

	return nil, nil
}

// WEB REST API used import data to system
func restImpexImport(params *api.T_APIHandlerParams) (interface{}, error) {

	modelName := params.RequestURLParams["model"]
	model, err := models.GetModel(modelName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	object, isObject := model.(models.I_Object)
	_, isStorable := model.(models.I_Storable)
	if !isObject || !isStorable {
		return nil, env.ErrorNew(modelName + " not implements I_Object or I_Storable interface")
	}

	attributes := make(map[string]models.T_AttributeInfo)
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
			return nil, env.ErrorNew("there is no attribute " + csvColumn)
		}
	}

	for csvRecord, err := csvReader.Read(); err == nil; csvRecord, err = csvReader.Read() {
		model, _ = model.New()
		object, _ = model.(models.I_Object)
		for idx, value := range csvRecord {
			object.Set(csvColumns[idx], value)
		}
		object.(models.I_Storable).Save()
	}

	return nil, env.ErrorNew("not implemented")
}
