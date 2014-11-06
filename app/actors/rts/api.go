package rts

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func setupAPI() error {
	var err error

	// 1. DefaultRtsAPI
	//----------------------
	err = api.GetRestService().RegisterAPI("rts", "GET", "visit", restRegisterVisit)
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

	err = api.GetRestService().RegisterAPI("rts", "GET", "visitors/realtime", restGetVisitsRealtime)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

func restRegisterVisit(params *api.T_APIHandlerParams) (interface{}, error) {
	eventData := make(map[string]interface{})
	session := params.Session
	eventData["sessionId"] = session.GetId()

	env.Event("api.visits", eventData)

	eventData = make(map[string]interface{})
	xReferrer := utils.InterfaceToString(params.Request.Header.Get("X-Referer"))

	eventData["referrer"] = xReferrer
	eventData["sessionId"] = session.GetId()

	http.SetCookie(params.ResponseWriter, &http.Cookie{Name: "X_Referrer", Value: xReferrer, Path: "/"})

	env.Event("api.referrer", eventData)
	env.Event("api.regVisitorAsOnlineHandler", eventData)

	return nil, nil
}

func restGetReferrers(params *api.T_APIHandlerParams) (interface{}, error) {
	result := make(map[string]int)

	for url := range referrers {
		result[utils.InterfaceToString(url)] = referrers[utils.InterfaceToString(url)]
	}

	return result, nil
}

func restGetVisits(params *api.T_APIHandlerParams) (interface{}, error) {
	result := make(map[string]interface{})

	err := GetTodayVisitorsData()
	if err != nil {
		return result, nil
	}

	result["visitsToday"] = visitorsInfoToday.Visitors
	result["ratio"] = 1

	err = GetYesterdayVisitorsData()
	if err != nil {
		return result, nil
	}
	countYesterday := visitorsInfoYesterday.Visitors
	countToday := visitorsInfoToday.Visitors
	if countYesterday != 0 {
		ratio := float64(countToday)/float64(countYesterday) - float64(1)
		result["ratio"] = utils.Round(ratio, 0.5, 2)
	}

	return result, nil
}

func restGetVisitsDetails(params *api.T_APIHandlerParams) (interface{}, error) {
	result := make(map[string]int)

	fromDatetmp, present := params.RequestURLParams["from"]
	if !present {
		fromDatetmp = time.Now().Format("2006-01-02")
	}
	fromDate, _ := time.Parse("2006-01-02", fromDatetmp)

	toDatetmp, present := params.RequestURLParams["to"]
	if !present {
		toDatetmp = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}
	toDate, _ := time.Parse("2006-01-02", toDatetmp)

	delta := toDate.Sub(fromDate)
	visitorInfoCollection, err := db.GetCollection(CollectionNameVisitors)
	if err == nil {
		visitorInfoCollection.AddFilter("day", ">=", fromDate)
		visitorInfoCollection.AddFilter("day", "<", toDate)
		visitorInfoCollection.AddSort("day", false)
		dbRecord, _ := visitorInfoCollection.Load()
		dbResult := make(map[string]int)
		if delta.Hours() > 48 {
			if len(dbRecord) > 0 {
				for _, item := range dbRecord {
					timestamp := fmt.Sprintf("%v", int32(utils.InterfaceToTime(item["day"]).Unix()))
					dbResult[timestamp] = utils.InterfaceToInt(item["visitors"])
				}
			}
			// group by days
			for date := fromDate; int32(date.Unix()) < int32(toDate.Unix()); date = date.AddDate(0, 0, 1) {
				timestamp := fmt.Sprintf("%v", int32(date.Unix()))
				result[timestamp] = dbResult[timestamp]
			}

		} else {
			if len(dbRecord) > 0 {
				details := RtsDecodeDetails(utils.InterfaceToString(dbRecord[0]["details"]))
				for _, item := range details {
					timestamp := fmt.Sprintf("%v", int32(utils.InterfaceToTime(item.Time).Unix()))
					dbResult[timestamp]++
				}
			}
			//	group by hours
			for date := fromDate; int32(date.Unix()) < int32(toDate.Unix()); date = date.AddDate(0, 0, 1) {

				timestamp := int32(date.Unix())
				year := time.Unix(int64(timestamp), 0).Year()
				month := time.Unix(int64(timestamp), 0).Month()
				day := time.Unix(int64(timestamp), 0).Day()
				for hour := 0; hour < 24; hour++ {
					timeGroup := time.Date(year, month, day, hour, 0, 0, 0, time.Local)
					if timeGroup.Unix() > time.Now().Unix() {
						break
					}
					timestamp := fmt.Sprintf("%v", int32(timeGroup.Unix()))
					result[timestamp] = dbResult[timestamp]
				}
			}
		}
	}

	return result, nil
}

func restGetConversions(params *api.T_APIHandlerParams) (interface{}, error) {
	result := make(map[string]interface{})

	result["totalVisitors"] = visitorsInfoToday.Visitors
	result["addedToCart"] = visitorsInfoToday.Cart
	result["reachedCheckout"] = visitorsInfoToday.Checkout
	result["purchased"] = visitorsInfoToday.Sales

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
	fromDatetmp, present := params.RequestURLParams["from"]
	if !present {
		fromDatetmp = time.Now().Format("2006-01-02")
	}
	fromDate, _ := time.Parse("2006-01-02", fromDatetmp)
	// check request params
	//---------------------
	toDatetmp, present := params.RequestURLParams["to"]
	if !present {
		toDatetmp = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}
	toDate, _ := time.Parse("2006-01-02", toDatetmp)

	h := md5.New()
	io.WriteString(h, fromDatetmp+"/"+toDatetmp)
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

	salesCollection, err := db.GetCollection(CollectionNameSales)
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
		productID := utils.InterfaceToString(item["product_id"])
		result[productID] = &SellerInfo{}
		if _, ok := topSellers.Data[productID]; !ok {

			product, err := product.LoadProductById(productID)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			productModel.SetId(productID)
			mediaPath, err := productModel.GetMediaPath("image")
			if err != nil {
				return result, env.ErrorDispatch(err)
			}

			if product.GetDefaultImage() != "" {
				result[productID].Image = mediaPath + product.GetDefaultImage()
			}

			result[productID].Name = product.GetName()
		}

		result[productID].Count = utils.InterfaceToInt(item["count"])
	}

	return result, nil
}

func restGetVisitsRealtime(params *api.T_APIHandlerParams) (interface{}, error) {
	result := make(map[string]interface{})
	ratio := float64(0)

	result["Online"] = len(OnlineSessions)
	if OnlineSessionsMax == 0 || len(OnlineSessions) == 0 {
		ratio = float64(1)
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

//"6CMAbByF0NXLo3SYnhKewcVU3QvTBIV0":{"Time":"2014-11-03T14:00:00+02:00","Checkout":0},
//"8q8prRhwX867xA2O815wwfkbA3O7UmYu":{"Time":"2014-11-03T14:00:00+02:00","Checkout":0},
//"E5LJgFNDfRnAzK01xS2rZCKJGGo8kIef":{"Time":"2014-11-03T15:00:00+02:00","Checkout":0},
//"IigFNY9wBLdREeuk9RiLomZKeTeFLrEj":{"Time":"2014-11-03T13:00:00+02:00","Checkout":0},
//"JzMTe8gUKmNTs7cqyTaJPUTKgmAYqCdn":{"Time":"2014-11-03T15:00:00+02:00","Checkout":0},
//"P4jSIKvyAy8AuE5jM6xVPZuKgnKm6nhW":{"Time":"2014-11-03T12:00:00+02:00","Checkout":0},
//"QYkn0SOp7OSoHjdOk9nGr8Kf3ZIexQLG":{"Time":"2014-11-03T14:00:00+02:00","Checkout":0},
//"S6JWtufTE3CJ9pSqWmRaOz6ZDZCKZ5gn":{"Time":"2014-11-03T13:00:00+02:00","Checkout":0},
//"XRGjloh5WCfBdk3sRHWTOelvN9mRSHSa":{"Time":"2014-11-03T15:00:00+02:00","Checkout":0},
//"XpOl2skeLC2cpKwrf0IvVZ6VfKyg2UH7":{"Time":"2014-11-03T13:00:00+02:00","Checkout":0},
//"ZGpzT2YBIGxLv8EISknOpWSr9F0OofFG":{"Time":"2014-11-03T14:00:00+02:00","Checkout":3},
//"aHcAyXOu5sdQcSt14820GktSKbRoI5BT":{"Time":"2014-11-03T14:00:00+02:00","Checkout":0},
//"axqKlgPDPDcJVWrA1VvCxWarzIqvECSs":{"Time":"2014-11-03T13:00:00+02:00","Checkout":0},
//"bGXau8aXIxoy54hUSDCqHhywsBBe9VCn":{"Time":"2014-11-03T14:00:00+02:00","Checkout":0},
//"bsALDa1sidVlTbaBarax6TuCFl73gNtJ":{"Time":"2014-11-03T14:00:00+02:00","Checkout":0},
//"caI4sH7w3Ucs2zFoQ28V0TkuIGVecjN0":{"Time":"2014-11-03T13:00:00+02:00","Checkout":0},
//"eOXZFY8DmFQuBCsLM95z2uZYvB1q7ToC":{"Time":"2014-11-03T13:00:00+02:00","Checkout":0},
//"hiT0ZSS5FDOYPf7k7YPPHo8N3RFNJBYK":{"Time":"2014-11-03T13:00:00+02:00","Checkout":0},
//"j7qZjO9QFQKUo1oWX0puYaQYVBh8DWSH":{"Time":"2014-11-03T15:00:00+02:00","Checkout":0},
//"sqQzO5w9D6ABnAFixZYIBlf0qrgmUPUK":{"Time":"2014-11-03T12:00:00+02:00","Checkout":0},
//"xnOua06f3j9atnJa8BVDGEWdEMpaiFtk":{"Time":"2014-11-03T14:00:00+02:00","Checkout":0}}
