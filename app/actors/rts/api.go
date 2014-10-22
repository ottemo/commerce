package rts

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"time"
	"fmt"
)

func setupAPI() error {
	var err error = nil

	// 1. DefaultProduct API
	//----------------------
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

	return nil
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

	if 0 != len(visits.Data[visits.Yesterday]) {
		countYesterday := len(visits.Data[visits.Yesterday])
		countToday := len(visits.Data[visits.Today])
		ratio := float64(countToday) / float64(countYesterday) - float64(1)
		result["ratio"] = utils.Round(ratio, 0.5, 2)
	} else {
		result["ratio"] = 100
	}

	return result, nil
}

func restGetVisitsDetails(params *api.T_APIHandlerParams) (interface{}, error) {
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

	delta := toDate.Sub(fromDate)

	if delta.Hours() > 48 {
		// group by days
		for date := fromDate; int32(date.Unix()) < int32(toDate.Unix()); date = date.AddDate(0, 0, 1) {
			timestamp := fmt.Sprintf("%v", int32(date.Unix()))
			result[timestamp] = len(visits.Data[date.Format("2006-01-02")])
		}
	} else {
		// group by days
		for date := fromDate; int32(date.Unix()) < int32(toDate.Unix()); date = date.AddDate(0, 0, 1) {
			for _, timestamp := range visits.Data[date.Format("2006-01-02")] {
				hour := time.Unix(int64(timestamp), 0).Hour()
				year := time.Unix(int64(timestamp), 0).Year()
				month := time.Unix(int64(timestamp), 0).Month()
				day := time.Unix(int64(timestamp), 0).Day()
				timeGroup := time.Date(year, month, day, hour, 0, 0, 0, time.Local)

				mapIndex := fmt.Sprintf("%v", int32(timeGroup.Unix()))
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
