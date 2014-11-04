package rts

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/utils"
	"strings"
	"time"
)

func referrerHandler(event string, data map[string]interface{}) bool {

	if "api.referrer" != event || "" == utils.InterfaceToString(data["referrer"]) {
		return true
	}

	referrer, err := GetReferrer(utils.InterfaceToString(data["referrer"]))
	if err != nil {
		return true
	}

	referrers[referrer] += 1

	return true
}

func visitsHandler(event string, data map[string]interface{}) bool {

	if "api.visits" != event {
		return true
	}

	err := GetTodayVisitorsData()
	if err != nil {
		return true
	}
	sessionId := utils.InterfaceToString(data["sessionId"])

	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	hour := time.Now().Hour()
	today := time.Date(year, month, day, hour, 0, 0, 0, time.Local)

	if _, ok := visitorsInfoToday.Details[sessionId]; !ok {
		visitorsInfoToday.Details[sessionId] = &VisitorDetail{Time: today}
		visitorsInfoToday.Visitors += 1
	}

	visitorsInfoToday.Details[sessionId] = &VisitorDetail{Time: today}
	_ = SaveVisitorData()

	return true
}

func addToCartHandler(event string, data map[string]interface{}) bool {

	if "api.addToCart" != event {
		return true
	}

	err := GetTodayVisitorsData()
	if err != nil {
		return true
	}
	sessionId := utils.InterfaceToString(data["sessionId"])

	if 0 == visitorsInfoToday.Details[sessionId].Checkout {
		visitorsInfoToday.Details[sessionId].Checkout = VISITOR_ADD_TO_CART
		visitorsInfoToday.Cart += 1
	}

	_ = SaveVisitorData()

	return true
}

func reachedCheckoutHandler(event string, data map[string]interface{}) bool {

	if "api.reachedCheckout" != event {
		return true
	}

	err := GetTodayVisitorsData()
	if err != nil {
		return true
	}
	sessionId := utils.InterfaceToString(data["sessionId"])

	if VISITOR_CHECKOUT > visitorsInfoToday.Details[sessionId].Checkout {
		visitorsInfoToday.Details[sessionId].Checkout = VISITOR_CHECKOUT
		visitorsInfoToday.Checkout += 1
	}

	_ = SaveVisitorData()

	return true
}

func purchasedHandler(event string, data map[string]interface{}) bool {

	if "api.purchased" != event {
		return true
	}

	err := GetTodayVisitorsData()
	if err != nil {
		return true
	}
	sessionId := utils.InterfaceToString(data["sessionId"])

	if VISITOR_SALES > visitorsInfoToday.Details[sessionId].Checkout {
		visitorsInfoToday.Details[sessionId].Checkout = VISITOR_SALES
		visitorsInfoToday.Sales += 1
	}

	_ = SaveVisitorData()

	return true
}

func salesHandler(event string, data map[string]interface{}) bool {

	if "api.sales" != event || len(data) == 0 {
		return true
	}
	salesData := make(map[string]int)

	salesHistoryCollection, err := db.GetCollection(COLLECTION_NAME_SALES_HISTORY)
	if err != nil {
		return true
	}

	for productId, count := range data {
		year := time.Now().Year()
		month := time.Now().Month()
		day := time.Now().Day()
		date := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		salesHistoryRow := make(map[string]interface{})
		salesData[productId] = utils.InterfaceToInt(count)

		salesHistoryCollection.ClearFilters()
		salesHistoryCollection.AddFilter("created_at", "=", date)
		salesHistoryCollection.AddFilter("product_id", "=", productId)
		dbSaleRow, _ := salesHistoryCollection.Load()

		newCount := utils.InterfaceToInt(count)
		if len(dbSaleRow) > 0 {
			salesHistoryRow["_id"] = utils.InterfaceToString(dbSaleRow[0]["_id"])
			oldCount := utils.InterfaceToInt(dbSaleRow[0]["count"])
			newCount += oldCount
		}

		// Add history row
		salesHistoryRow["product_id"] = productId
		salesHistoryRow["created_at"] = date
		salesHistoryRow["count"] = newCount
		_, err = salesHistoryCollection.Save(salesHistoryRow)
		if err != nil {
			return true
		}
	}

	SaveSalesData(salesData)

	return true
}

func regVisitorAsOnlineHandler(event string, data map[string]interface{}) bool {
	if "api.regVisitorAsOnlineHandler" != event {
		return true
	}
	sessionId := utils.InterfaceToString(data["sessionId"])

	referrerType := REFERRER_TYPE_DIRECT

	if "" != utils.InterfaceToString(data["referrer"]) {
		referrer, err := GetReferrer(utils.InterfaceToString(data["referrer"]))
		if err != nil {
			return true
		}

		isSearchEngine := false
		for index := 0; index < len(searchEngines); index += 1 {
			if strings.Contains(referrer, searchEngines[index]) {
				isSearchEngine = true
			}
		}

		if isSearchEngine {
			referrerType = REFERRER_TYPE_SEARCH
		} else {
			referrerType = REFERRER_TYPE_SITE
		}
	}

	if _, ok := OnlineSessions[sessionId]; !ok {
		OnlineSessions[sessionId] = &OnlineReferrer{}
		IncreaseOnline(referrerType)
		if len(OnlineSessions) > OnlineSessionsMax {
			OnlineSessionsMax = len(OnlineSessions)
		}
	} else {
		if OnlineSessions[sessionId].referrerType != referrerType {
			DecreaseOnline(OnlineSessions[sessionId].referrerType)
			IncreaseOnline(referrerType)
		}
	}


	OnlineSessions[sessionId].time = time.Now()
	OnlineSessions[sessionId].referrerType = referrerType

	return true
}

func visitorOnlineActionHandler(event string, data map[string]interface{}) bool {
	if "api.visitorOnlineAction" != event {
		return true
	}
	sessionId := utils.InterfaceToString(data["sessionId"])
	if _, ok := OnlineSessions[sessionId]; ok {
		OnlineSessions[sessionId].time = time.Now()
	}

	return true
}
