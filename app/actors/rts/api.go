package rts

import (
	"crypto/md5"
	"fmt"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"io"
	"time"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/app/models/product"
)

func setupAPI() error {
	var err error = nil

	// 1. DefaultRtsAPI
	//----------------------
	err = api.GetRestService().RegisterAPI("rts", "GET", "visit", restRegVisit)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts", "GET", "referrers", restGetReferrers)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts", "GET", "visits", restGetVisits)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts", "GET", "visits/details/:from/:to", restGetVisitsDetails)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts", "GET", "conversions", restGetConversions)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts", "GET", "sales", restGetSales)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts", "GET", "sales/details/:from/:to", restGetSalesDetails)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts", "GET", "top_sellers", restGetTopSellers)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

func restRegVisit(params *api.T_APIHandlerParams) (interface{}, error) {
	eventData := make(map[string]interface{})
	session := params.Session
	eventData["sessionId"] = session.GetId()

	env.Event("api.visits", eventData)

	return nil, nil
}

func restGetReferrers(params *api.T_APIHandlerParams) (interface{}, error) {
	result := make(map[string]int)

	for url, _ := range referrers {
		result[utils.InterfaceToString(url)] = referrers[utils.InterfaceToString(url)].Count
	}

	return result, nil
}

func restGetVisits(params *api.T_APIHandlerParams) (interface{}, error) {
	result := make(map[string]interface{})

	result["visitsToday"] = len(visits.Data[visits.Today])
	result["ratio"] = 0

	if 0 != len(visits.Data[visits.Yesterday]) {
		countYesterday := len(visits.Data[visits.Yesterday])
		countToday := len(visits.Data[visits.Today])
		ratio := float64(countToday) / float64(countYesterday) - float64(1)
		result["ratio"] = utils.Round(ratio, 0.5, 2)
	}

	return result, nil
}

func restGetVisitsDetails(params *api.T_APIHandlerParams) (interface{}, error) {
	result := make(map[string]int)

	fromDate_tmp, present := params.RequestURLParams["from"]
	if !present {
		fromDate_tmp = time.Now().Format("2006-01-02")
	}
	fromDate, _ := time.Parse("2006-01-02", fromDate_tmp)

	toDate_tmp, present := params.RequestURLParams["to"]
	if !present {
		toDate_tmp = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}
	toDate, _ := time.Parse("2006-01-02", toDate_tmp)

	delta := toDate.Sub(fromDate)

	if delta.Hours() > 48 {
		// group by days
		for date := fromDate; int32(date.Unix()) < int32(toDate.Unix()); date = date.AddDate(0, 0, 1) {
			timestamp := fmt.Sprintf("%v", int32(date.Unix()))
			result[timestamp] = len(visits.Data[date.Format("2006-01-02")])
		}
	} else {
		// group by hours
		for date := fromDate; int32(date.Unix()) < int32(toDate.Unix()); date = date.AddDate(0, 0, 1) {

			timestamp := int32(date.Unix())
			year := time.Unix(int64(timestamp), 0).Year()
			month := time.Unix(int64(timestamp), 0).Month()
			day := time.Unix(int64(timestamp), 0).Day()
			for hour := 0; hour < 24; hour += 1 {
				timeGroup := time.Date(year, month, day, hour, 0, 0, 0, time.Local)
				if timeGroup.Unix() > time.Now().Unix() {
					break
				}
				mapIndex := fmt.Sprintf("%v", int32(timeGroup.Unix()))
				result[mapIndex] = 0
			}
			for _, timestamp := range visits.Data[date.Format("2006-01-02")] {
				mapIndex := GetDayForTimestamp(int64(timestamp), true)
				result[mapIndex] += 1
			}
		}
	}

	return result, nil
}

func restGetConversions(params *api.T_APIHandlerParams) (interface{}, error) {
	result := make(map[string]interface{})

	result["totalVisitors"] = conversions["visitors"]["count"]
	result["addedToCart"] = len(conversions["addedToCart"])
	result["reachedCheckout"] = len(conversions["reachedCheckout"])
	result["purchased"] = len(conversions["purchased"])

	return result, nil
}

func restGetSales(params *api.T_APIHandlerParams) (interface{}, error) {
	result := make(map[string]interface{})

	currDate, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	yesterdayDate, _ := time.Parse("2006-01-02", currDate.AddDate(0, 0, -1).Format("2006-01-02"))

	if sales.lastUpdate == 0 { // Init sales data
		GetTotalSales(currDate, yesterdayDate)
		result["today"] = sales.today
		result["ratio"] = sales.ratio
	} else {
		lastUpdate, _ := time.Parse("2006-01-02", time.Unix(int64(sales.lastUpdate), 0).Format("2006-01-02"))
		delta := currDate.Sub(lastUpdate)
		if delta > 1 { // Updates the sales data if they older 1 hour
			GetTotalSales(currDate, yesterdayDate)
			result["today"] = sales.today
			result["ratio"] = sales.ratio
		} else {
			// Returns the  existing data
			result["today"] = sales.today
			result["ratio"] = sales.ratio
		}
	}

	return result, nil
}

func restGetSalesDetails(params *api.T_APIHandlerParams) (interface{}, error) {
	result := make(map[string]int)

	// check request params
	//---------------------
	fromDate_tmp, present := params.RequestURLParams["from"]
	if !present {
		fromDate_tmp = time.Now().Format("2006-01-02")
	}
	fromDate, _ := time.Parse("2006-01-02", fromDate_tmp)
	// check request params
	//---------------------
	toDate_tmp, present := params.RequestURLParams["to"]
	if !present {
		toDate_tmp = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}
	toDate, _ := time.Parse("2006-01-02", toDate_tmp)

	h := md5.New()
	io.WriteString(h, fromDate_tmp+"/"+toDate_tmp)
	periodHash := fmt.Sprintf("%x", h.Sum(nil))

	if _, ok := salesDetail[periodHash]; !ok {
		salesDetail[periodHash] = &SalesDetailData{Data: make(map[string]int)}

		GetSalesDetail(fromDate, toDate, periodHash)

	} else {
		// check last updates
		if salesDetail[periodHash].lastUpdate == 0 {

			GetSalesDetail(fromDate, toDate, periodHash)

		} else {
			currDate, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
			lastUpdate, _ := time.Parse("2006-01-02", time.Unix(int64(salesDetail[periodHash].lastUpdate), 0).Format("2006-01-02"))
			delta := currDate.Sub(lastUpdate)
			if delta > 1 { // Updates the sales data if they older than 1 hour
				GetSalesDetail(fromDate, toDate, periodHash)
			}
		}
	}

	result = salesDetail[periodHash].Data

	return result, nil
}

func restGetTopSellers(params *api.T_APIHandlerParams) (interface{}, error) {
	result := make(map[string]*SellerInfo)

	salesCollection, err := db.GetCollection(COLLECTION_NAME_SALES)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	salesCollection.AddFilter("count", ">", 0)
	salesCollection.AddSort("count", true)
	salesCollection.SetLimit(0, 5)
	items, _ := salesCollection.Load()

	productModel, err := product.GetProductModel()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, item := range items {
		productId := utils.InterfaceToString(item["product_id"])
		result[productId] = &SellerInfo{}
		if _, ok := topSellers.Data[productId]; !ok {

			product, err := product.LoadProductById(productId)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			productModel.SetId(productId)
			mediaPath, err := productModel.GetMediaPath("image")
			if err != nil {
				return result, env.ErrorDispatch(err)
			}

			if product.GetDefaultImage() != "" {
				result[productId].Image = mediaPath + product.GetDefaultImage()
			}

			result[productId].Name = product.GetName()
		}

		result[productId].Count = utils.InterfaceToInt(item["count"])
	}

	return result, nil
}
