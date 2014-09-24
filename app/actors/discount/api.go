package discount

import (
	"encoding/csv"
	"errors"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/utils"
)

// initializes API for discount
func setupAPI() error {
	var err error = nil

	err = api.GetRestService().RegisterAPI("discount", "GET", "apply/:code", restDiscountApply)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("discount", "GET", "neglect/:code", restDiscountNeglect)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("discount", "GET", "download/csv", restDiscountCSVDownload)
	if err != nil {
		return err
	}

	err = api.GetRestService().RegisterAPI("discount", "POST", "upload/csv", restDiscountCSVUpload)
	if err != nil {
		return err
	}

	return nil
}

// WEB REST API function to apply discount code to current checkout
func restDiscountApply(params *api.T_APIHandlerParams) (interface{}, error) {

	couponCode := params.RequestURLParams["code"]

	// getting applied coupons array for current session
	appliedCoupons := make([]string, 0)
	if sessionValue, ok := params.Session.Get(SESSION_KEY_APPLIED_DISCOUNT_CODES).([]string); ok {
		appliedCoupons = sessionValue
	}

	// checking if coupon was already applied
	if utils.IsInArray(couponCode, appliedCoupons) {
		return nil, errors.New("coupon code already applied")
	}

	// loading coupon for specified code
	collection, err := db.GetCollection(COLLECTION_NAME_COUPON_DISCOUNTS)
	if err != nil {
		return nil, err
	}
	err = collection.AddFilter("code", "=", couponCode)
	if err != nil {
		return nil, err
	}

	records, err := collection.Load()
	if err != nil {
		return nil, err
	}

	// checking and applying obtained coupon
	if len(records) > 0 {
		discountCoupon := records[0]

		applyTimes := utils.InterfaceToInt(discountCoupon["times"])
		workSince := utils.InterfaceToTime(discountCoupon["since"])
		workUntil := utils.InterfaceToTime(discountCoupon["until"])

		currentTime := time.Now()

		// to be applicable coupon should satisfy following conditions:
		//   [applyTimes] should be -1 or >0 and [workSince] >= currentTime <= [workUntil] if set
		if (applyTimes == -1 || applyTimes > 0) &&
			(utils.IsZeroTime(workSince) || workSince.Unix() <= currentTime.Unix()) &&
			(utils.IsZeroTime(workUntil) || workUntil.Unix() >= currentTime.Unix()) {

			// TODO: coupon loosing with session clear, probably should be made on order creation, or have event on session
			// times used decrease
			if applyTimes > 0 {
				discountCoupon["times"] = applyTimes - 1
				_, err := collection.Save(discountCoupon)
				if err != nil {
					return nil, err
				}
			}

			// coupon is working - applying it
			appliedCoupons = append(appliedCoupons, couponCode)
			params.Session.Set(SESSION_KEY_APPLIED_DISCOUNT_CODES, appliedCoupons)

		} else {
			return nil, errors.New("coupon is not applicable")
		}
	} else {
		return nil, errors.New("coupon code not found")
	}

	return "ok", nil
}

// WEB REST API function to neglect(un-apply) discount code to current checkout
//   - use "*" as code to neglect all discounts
func restDiscountNeglect(params *api.T_APIHandlerParams) (interface{}, error) {

	couponCode := params.RequestURLParams["code"]

	if couponCode == "*" {
		params.Session.Set(SESSION_KEY_APPLIED_DISCOUNT_CODES, make([]string, 0))
		return "ok", nil
	}

	if appliedCoupons, ok := params.Session.Get(SESSION_KEY_APPLIED_DISCOUNT_CODES).([]string); ok {
		newAppliedCoupons := make([]string, 0)
		for _, value := range appliedCoupons {
			if value != couponCode {
				newAppliedCoupons = append(newAppliedCoupons, value)
			}
		}
		params.Session.Set(SESSION_KEY_APPLIED_DISCOUNT_CODES, newAppliedCoupons)

		// times used increase
		collection, err := db.GetCollection(COLLECTION_NAME_COUPON_DISCOUNTS)
		if err != nil {
			return nil, err
		}
		err = collection.AddFilter("code", "=", couponCode)
		if err != nil {
			return nil, err
		}
		records, err := collection.Load()
		if err != nil {
			return nil, err
		}
		if len(records) > 0 {
			applyTimes := utils.InterfaceToInt(records[0]["times"])
			if applyTimes >= 0 {
				records[0]["times"] = applyTimes + 1

				_, err := collection.Save(records[0])
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return "ok", nil
}

// WEB REST API function to download current tax rates in CSV format
func restDiscountCSVDownload(params *api.T_APIHandlerParams) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, err
	}

	// preparing csv writer
	csvWriter := csv.NewWriter(params.ResponseWriter)

	params.ResponseWriter.Header().Set("Content-type", "text/csv")
	params.ResponseWriter.Header().Set("Content-disposition", "attachment;filename=discount_coupons.csv")

	csvWriter.Write([]string{"Code", "Name", "Amount", "Percent", "Times", "Since", "Until"})
	csvWriter.Flush()

	// loading records from DB and writing them in csv format
	collection, err := db.GetCollection(COLLECTION_NAME_COUPON_DISCOUNTS)
	if err != nil {
		return nil, err
	}

	err = collection.Iterate(func(record map[string]interface{}) bool {
		csvWriter.Write([]string{
			utils.InterfaceToString(record["code"]),
			utils.InterfaceToString(record["name"]),
			utils.InterfaceToString(record["amount"]),
			utils.InterfaceToString(record["percent"]),
			utils.InterfaceToString(record["times"]),
			utils.InterfaceToString(record["since"]),
			utils.InterfaceToString(record["until"]),
		})

		csvWriter.Flush()
		return true
	})

	return nil, nil
}

// WEB REST API function to upload tax rates into CSV
func restDiscountCSVUpload(params *api.T_APIHandlerParams) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, err
	}

	csvFile, _, err := params.Request.FormFile("file")
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ','

	collection, err := db.GetCollection(COLLECTION_NAME_COUPON_DISCOUNTS)
	if err != nil {
		return nil, err
	}
	collection.Delete()

	csvReader.Read() //skipping header
	for csvRecord, err := csvReader.Read(); err == nil; csvRecord, err = csvReader.Read() {
		if len(csvRecord) >= 7 {
			record := make(map[string]interface{})

			code := utils.InterfaceToString(csvRecord[0])
			name := utils.InterfaceToString(csvRecord[1])
			if code == "" || name == "" {
				continue
			}

			times := utils.InterfaceToInt(csvRecord[4])
			if csvRecord[4] == "" {
				times = -1
			}

			record["code"] = code
			record["name"] = name
			record["amount"] = utils.InterfaceToFloat64(csvRecord[2])
			record["percent"] = utils.InterfaceToFloat64(csvRecord[3])
			record["times"] = times
			record["since"] = utils.InterfaceToTime(csvRecord[5])
			record["until"] = utils.InterfaceToTime(csvRecord[6])

			collection.Save(record)
		}
	}

	return "ok", nil
}
