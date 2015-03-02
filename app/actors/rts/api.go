package rts

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	var err error

	err = api.GetRestService().RegisterAPI("rts/visit", api.ConstRESTOperationCreate, APIRegisterVisit)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts/referrers", api.ConstRESTOperationGet, APIGetReferrers)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts/visits", api.ConstRESTOperationGet, APIGetVisits)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts/visits/detail/:from/:to", api.ConstRESTOperationGet, APIGetVisitsDetails)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts/conversion", api.ConstRESTOperationGet, APIGetConversion)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts/sales", api.ConstRESTOperationGet, APIGetSales)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts/sales/detail/:from/:to", api.ConstRESTOperationGet, APIGetSalesDetails)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts/bestsellers", api.ConstRESTOperationGet, APIGetBestsellers)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("rts/visits/realtime", api.ConstRESTOperationGet, APIGetVisitsRealtime)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APIRegisterVisit registers request for a statistics
func APIRegisterVisit(context api.InterfaceApplicationContext) (interface{}, error) {
	if httpRequest, ok := context.GetRequest().(*http.Request); ok && httpRequest != nil {
		if httpResponseWriter, ok := context.GetResponseWriter().(http.ResponseWriter); ok && httpResponseWriter != nil {
			xReferrer := utils.InterfaceToString(httpRequest.Header.Get("X-Referer"))

			http.SetCookie(httpResponseWriter, &http.Cookie{Name: "X_Referrer", Value: xReferrer, Path: "/"})

			eventData := map[string]interface{}{"session": context.GetSession(), "context": context}
			env.Event("api.rts.visit", eventData)

			return nil, nil
		}
	}
	return nil, nil
}

// APIGetReferrers returns list of unique referrers were registered
func APIGetReferrers(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]int)

	for url, count := range referrers {
		result[url] = count
	}

	return result, nil
}

// APIGetVisits returns site visit information for current day
func APIGetVisits(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]interface{})

	err := GetTodayVisitorsData()
	if err != nil {
		return result, nil
	}

	err = GetYesterdayVisitorsData()
	if err != nil {
		return result, nil
	}

	countToday := visitorsInfoToday.Visitors
	countYesterday := visitorsInfoYesterday.Visitors

	result["visitsToday"] = countToday
	result["visitsYesterday"] = countYesterday
	result["ratio"] = 1

	if countYesterday > 0 {
		ratio := float64(countToday)/float64(countYesterday) - float64(1)
		result["ratio"] = utils.Round(ratio, 0.5, 2)
	}

	return result, nil
}

// APIGetVisits returns site visit information for a day _ver 2
func APIGetVisitsVerTwo(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]interface{})
	timeZone := context.GetRequestArgument("tz")

	// get a hours pasted for local day
	todayTo := time.Now().Truncate(time.Hour)
	todayFrom, yesterdayFrom := GetLocalOneDayBefore(todayTo, timeZone)

	// get data for visits
	todayVisits, err := GetRangeVisits(todayFrom, todayTo)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	yesterdayVisits, _ := GetRangeVisits(yesterdayFrom, todayFrom)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// count ratio for current data
	ratio := float64(1)
	if yesterdayVisits !=0 {
		ratio = float64(todayVisits)/float64(yesterdayVisits) - float64(1)
	}

	// Returns the data
	result["visitsToday"] = todayVisits
	result["visitsYesterday"] = yesterdayVisits
	result["ratio"] = ratio

	return result, nil
}

// APIGetVisitsDetails returns detailed site visit information for a specified period
//   - period start and end dates should be specified in "from" and "to" attributes in DD-MM-YYY format
func APIGetVisitsDetails(context api.InterfaceApplicationContext) (interface{}, error) {

	// getting initial values
	result := make(map[string]int)
	timeZone := context.GetRequestArgument("tz")
	dateFrom := utils.InterfaceToTime(context.GetRequestArgument("from"))
	dateTo := utils.InterfaceToTime(context.GetRequestArgument("to"))

	// checking if user specified correct from and to dates
	if dateFrom.IsZero() {
		dateFrom = time.Now()
	}

	if dateTo.IsZero() {
		dateTo = time.Now()
	}

	if dateFrom == dateTo {
		dateTo = dateTo.Add(time.Hour * 24)
	}

	// time zone recognize routines
	if timeZone != "" {
		dateFrom = utils.ApplyTimeZone(dateFrom, timeZone)
		dateTo = utils.ApplyTimeZone(dateTo, timeZone)
	}

	dateFrom = dateFrom.Truncate(time.Hour * 24)
	dateTo = dateTo.Truncate(time.Hour * 24)

	// determining required scope
	delta := dateTo.Sub(dateFrom)

	timeScope := time.Hour
	if delta.Hours() > 48 {
		timeScope = timeScope * 24
	}

	// making database request
	visitorInfoCollection, err := db.GetCollection(ConstCollectionNameRTSVisitors)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorInfoCollection.AddFilter("day", ">=", dateFrom)
	visitorInfoCollection.AddFilter("day", "<", dateTo)
	visitorInfoCollection.AddSort("day", false)

	dbRecords, err := visitorInfoCollection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// filling requested period
	timeIterator := dateFrom
	for timeIterator.Before(dateTo) {
		result[fmt.Sprint(timeIterator.Unix())] = 0
		timeIterator = timeIterator.Add(timeScope)
	}

	// grouping database records
	for _, item := range dbRecords {
		timestamp := fmt.Sprint(utils.InterfaceToTime(item["day"]).Truncate(timeScope).Unix())
		visits := utils.InterfaceToInt(item["visitors"])

		if value, present := result[timestamp]; present {
			result[timestamp] = value + visits
		}
	}

	return result, nil
}

// APIGetConversion returns site conversion information
func APIGetConversion(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]interface{})

	result["totalVisitors"] = visitorsInfoToday.Visitors
	result["addedToCart"] = visitorsInfoToday.Cart
	result["reachedCheckout"] = visitorsInfoToday.Checkout
	result["purchased"] = visitorsInfoToday.Sales

	return result, nil
}


// NEW type of get sales
func APIGetSales(context api.InterfaceApplicationContext) (interface{}, error) {

	result := make(map[string]interface{})
	timeZone := context.GetRequestArgument("tz")

	// get a hours pasted for local day
	todayTo := time.Now().Truncate(time.Hour)
	todayFrom, yesterdayFrom := GetLocalOneDayBefore(todayTo, timeZone)

	// get data for sales

	todaySales, err := GetRangeSales(todayFrom, todayTo)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	yesterdaySales, _ := GetRangeSales(yesterdayFrom, todayFrom)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// count ratio for current data
	ratio := float64(1)
	if yesterdaySales !=0 {
		ratio = float64(todaySales)/float64(yesterdaySales) - float64(1)
	}

	// Returns the data
	result["today"] = todaySales
	result["yesterday"] = yesterdaySales
	result["ratio"] = ratio

	return result, nil
}


// APIGetSalesDetails returns site sales information for a specified period
//   - period start and end dates should be specified in "from" and "to" attributes in DD-MM-YYY format
func APIGetSalesDetails(context api.InterfaceApplicationContext) (interface{}, error) {

	// getting initial values
	result := make(map[string]int)
	timeZone := context.GetRequestArgument("tz")
	dateFrom := utils.InterfaceToTime(context.GetRequestArgument("from"))
	dateTo := utils.InterfaceToTime(context.GetRequestArgument("to"))

	currentTime := time.Now()

	// checking if user specified correct from and to dates
	if dateFrom.IsZero() {
		dateFrom = currentTime
	}

	if dateTo.IsZero() {
		dateTo = currentTime
	}

	if dateFrom == dateTo {
		dateTo = dateTo.Add(time.Hour * 24)
	}

	// time zone recognize routines
	dateFrom = utils.ApplyTimeZone(dateFrom, timeZone)
	dateTo = utils.ApplyTimeZone(dateTo, timeZone)

	dateFrom = dateFrom.Truncate(time.Hour * 24)
	dateTo = dateTo.Truncate(time.Hour * 24)

	// determining required scope
	delta := dateTo.Sub(dateFrom)

	timeScope := time.Hour
	if delta.Hours() > 48 {
		timeScope = timeScope * 24
	}

	// set database request settings
	orderCollectionModelT, err := order.GetOrderCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	dbCollection := orderCollectionModelT.GetDBCollection()
	dbCollection.SetResultColumns("_id", "created_at")
	dbCollection.AddSort("created_at", false)
	dbCollection.AddFilter("created_at", ">=", dateFrom)
	dbCollection.AddFilter("created_at", "<=", dateTo)

	// get database records
	dbRecords, err := dbCollection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// filling requested period
	timeIterator := dateFrom
	for timeIterator.Before(dateTo) {
		result[fmt.Sprint(timeIterator.Unix())] = 0
		timeIterator = timeIterator.Add(timeScope)
	}

	// grouping database records
	for _, order := range dbRecords {
		timestamp := fmt.Sprint(utils.InterfaceToTime(order["created_at"]).Truncate(timeScope).Unix())

		if _, present := result[timestamp]; present {
			result[timestamp]++
		}
	}

	return result, nil
}

// APIGetBestsellers returns information on site bestsellers
func APIGetBestsellers(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]*SellerInfo)

	salesCollection, err := db.GetCollection(ConstCollectionNameRTSSales)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	salesCollection.AddFilter("count", ">", 0)
	salesCollection.AddSort("count", true)
	salesCollection.SetLimit(0, 5)
	collectionRecords, _ := salesCollection.Load()

	for _, item := range collectionRecords {
		productID := utils.InterfaceToString(item["product_id"])
		result[productID] = &SellerInfo{}

		productInstance, err := product.LoadProductByID(productID)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		mediaPath, err := productInstance.GetMediaPath("image")
		if err != nil {
			return result, env.ErrorDispatch(err)
		}

		if productInstance.GetDefaultImage() != "" {
			result[productID].Image = mediaPath + productInstance.GetDefaultImage()
		}

		result[productID].Name = productInstance.GetName()
		result[productID].Count = utils.InterfaceToInt(item["count"])
	}

	return result, nil
}

// APIGetVisitsRealtime returns real-time information on current visits
func APIGetVisitsRealtime(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]interface{})
	ratio := float64(0)

	result["Online"] = len(OnlineSessions)
	if OnlineSessionsMax == 0 || len(OnlineSessions) == 0 {
		ratio = float64(0)
	} else {
		ratio = float64(len(OnlineSessions)) / float64(OnlineSessionsMax)
	}
	result["OnlineRatio"] = utils.Round(ratio, 0.5, 2)

	result["Direct"] = OnlineDirect
	if OnlineDirectMax == 0 || OnlineDirect == 0 {
		ratio = float64(0)
	} else {
		ratio = float64(OnlineDirect) / float64(OnlineDirectMax)
	}
	result["DirectRatio"] = utils.Round(ratio, 0.5, 2)

	result["Search"] = OnlineSearch
	if OnlineSearchMax == 0 || OnlineSearch == 0 {
		ratio = float64(0)
	} else {
		ratio = float64(OnlineSearch) / float64(OnlineSearchMax)
	}
	result["SearchRatio"] = utils.Round(ratio, 0.5, 2)

	result["Site"] = OnlineSite
	if OnlineSiteMax == 0 || OnlineSite == 0 {
		ratio = float64(0)
	} else {
		ratio = float64(OnlineSite) / float64(OnlineSiteMax)
	}
	result["SiteRatio"] = utils.Round(ratio, 0.5, 2)

	return result, nil
}
