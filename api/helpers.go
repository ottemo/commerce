package api

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/ottemo/foundation/api/session"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// returns session instance by id or nil
func GetSessionById(sessionId string) InterfaceSession {
	session, _ := session.GetSessionById(sessionId)
	return session
}

// returns true if admin rights allowed for current session
func ValidateAdminRights(params *StructAPIHandlerParams) error {
	if value := params.Session.Get(ConstSessionKeyAdminRights); value != nil {
		if value.(bool) == true {
			return nil
		}
	}

	return env.ErrorNew("no admin rights")
}

// tries to represent HTTP request content in map[string]interface{} format
func GetRequestContentAsMap(params *StructAPIHandlerParams) (map[string]interface{}, error) {

	result, ok := params.RequestContent.(map[string]interface{})
	if !ok {
		if params.Request.Method == "POST" {
			return nil, env.ErrorNew("unexpected request content")
		} else {
			result = make(map[string]interface{})
		}
	}

	return result, nil
}

// modifies collection with applying filters from request URL
func ApplyFilters(params *StructAPIHandlerParams, collection db.InterfaceDBCollection) error {

	for attributeName, attributeValue := range params.RequestGETParams {
		switch attributeName {
		case "limit":
			collection.SetLimit(GetListLimit(params))
		case "sort":
			attributesList := strings.Split(attributeValue, ",")

			for _, attributeName := range attributesList {
				descOrder := false
				if attributeName[0] == '^' {
					descOrder = true
					attributeName = strings.Trim(attributeName, "^")
				}
				collection.AddSort(attributeName, descOrder)
			}

		default:
			if collection.HasColumn(attributeName) {

				filterOperator := "="
				for _, prefix := range []string{">=", "<=", "!=", ">", "<", "~"} {
					if strings.HasPrefix(attributeValue, prefix) {
						attributeValue = strings.TrimPrefix(attributeValue, prefix)
						filterOperator = prefix
					}
				}
				if filterOperator == "~" {
					filterOperator = "like"
				}

				switch {
				case strings.Contains(attributeValue, ".."):
					rangeValues := strings.Split(attributeValue, "..")
					if rangeValues[0] != "" {
						collection.AddFilter(attributeName, ">=", rangeValues[0])
					}
					if rangeValues[1] != "" {
						collection.AddFilter(attributeName, "<=", rangeValues[1])
					}

				case strings.Contains(attributeValue, ","):
					options := strings.Split(attributeValue, ",")
					if filterOperator == "=" {
						collection.AddFilter(attributeName, "in", options)
					} else {
						collection.SetupFilterGroup(attributeName, true, "")
						for _, optionValue := range options {
							collection.AddGroupFilter(attributeName, attributeName, filterOperator, optionValue)
						}
					}

				default:
					collection.AddFilter(attributeName, filterOperator, attributeValue)
				}
			}
		}
	}
	return nil
}

// returns (offset, limit, error) values based on request string value
//   "1,2" will return offset: 1, limit: 2, error: nil
//   "2" will return offset: 0, limit: 2, error: nil
//   "something wrong" will return offset: 0, limit: 0, error: [error msg]
func GetListLimit(params *StructAPIHandlerParams) (int, int) {
	limitValue := ""

	if value, isLimit := params.RequestURLParams["limit"]; isLimit {
		limitValue = value
	} else if value, isLimit := params.RequestGETParams["limit"]; isLimit {
		limitValue = value
	} else {
		contentMap, err := GetRequestContentAsMap(params)
		if err == nil {
			if value, isLimit := contentMap["limit"]; isLimit {
				if value, ok := value.(string); ok {
					limitValue = value
				}
			}
		}
	}
	limitValue, _ = url.QueryUnescape(limitValue)

	splitResult := strings.Split(limitValue, ",")
	if len(splitResult) > 1 {
		offset, err := strconv.Atoi(strings.TrimSpace(splitResult[0]))
		if err != nil {
			return 0, 0
		}

		limit, err := strconv.Atoi(strings.TrimSpace(splitResult[1]))
		if err != nil {
			return 0, 0
		}

		return offset, limit
	} else if len(splitResult) > 0 {
		limit, err := strconv.Atoi(strings.TrimSpace(splitResult[0]))
		if err != nil {
			return 0, 0
		}

		return 0, limit
	}

	return 0, 0
}
