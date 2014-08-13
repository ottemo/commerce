package tax

import (
	"errors"
	"encoding/csv"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/api"

	"github.com/ottemo/foundation/app/utils"
)



// initializes API for tax
func setupAPI() error {
	var err error = nil

	err = api.GetRestService().RegisterAPI("tax", "GET", "upload/csv", restTaxCSVUpload)
	if err != nil {
		return err
	}

	return nil
}



// upload tax rates into CSV
func restTaxCSVUpload(params *api.T_APIHandlerParams) (interface{}, error) {

	csvFile, _, err := params.Request.FormFile("file")
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ','


	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Taxes"); err == nil {
			collection.Delete()

			csvReader.Read() //skipping header
			for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
				if len(record) >= 5 {
					taxRecord := make(map[string]interface{})

					taxRecord["code"] = record[0]
					taxRecord["country"] = record[1]
					taxRecord["state"] = record[2]
					taxRecord["zip"] = record[3]
					taxRecord["rate"] = utils.InterfaceToFloat64(record[4])

					collection.Save(taxRecord)
				}
			}
		} else {
			return nil, err
		}
	} else {
		return nil, errors.New("can't get DB engine")
	}

	return "ok", nil
}
