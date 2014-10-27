package rts

/**
 * referrers = {
 * 		'url_1': {
 * 			'Data': {
 *				'1970-01-01': {sessionID_1: 1, sessionID_2: 1, .., sessionID_N: 1}
 *				'1970-01-02': {sessionID_1: 1, sessionID_2: 1, .., sessionID_N: 1}
 *				'1970-01-03': {sessionID_1: 1, sessionID_2: 1, .., sessionID_N: 1}
 *			},
 *			'Count': N
 * 		},
 * 		'url_2': {
 * 			'Data': {
 *				'1970-01-01': {sessionID_1: 1, sessionID_2: 1, .., sessionID_N: 1}
 *				'1970-01-02': {sessionID_1: 1, sessionID_2: 1, .., sessionID_N: 1}
 *				'1970-01-03': {sessionID_1: 1, sessionID_2: 1, .., sessionID_N: 1}
 *			},
 *			'Count': N
 * 		},
 * }
 *
 * visits = {
 *		'Yesterday': 'DAY_2',
 *		'Today': 'DAY_1',
 *		'Data' : {
 * 			'DAY_1': {sessionID_1: timestamp, sessionID_2: timestamp, .., sessionID_N: timestamp},
 *			'...': 	 {sessionID_1: timestamp, sessionID_2: timestamp, .., sessionID_N: timestamp},
 *			'DAY_N': {sessionID_1: timestamp, sessionID_2: timestamp, .., sessionID_N: timestamp},
 *		}
 * }
 *
 * conversions = {
 *		"addedToCart": {sessionID_1: true, sessionID_2: true, .., sessionID_N: true},
 *		"reachedCheckout": {sessionID_1: true, sessionID_2: true, .., sessionID_N: true},
 *		"purchased": {sessionID_1: true, sessionID_2: true, .., sessionID_N: true},
 *		"visitors": X
 * }
 *
 * sales = {
 *		"lastUpdate": timestamp,
 *		"today": x,
 *		"yesterday": y,
 *		"ratio": z,
 * }
 *
 * salesDetail = {
 *		"period(MD5(dateFrom/dateTo))": {
 *			"Data": {},
 *			"lastUpdate": timestamp
 *		}
 * }
 *
 */

var (
	referrers   = make(map[string]*ReferrerData)
	visits      = Visits{Data: make(map[string]map[string]int32)}
	conversions = make(map[string]map[string]int)
	sales 		= Sales{}
	salesDetail = make(map[string]*SalesDetailData)
)

type ReferrerData struct {
	Data    		map[string]map[string]bool
	Count 			int
}

type Visits struct {
	Data 			map[string]map[string]int32
	Yesterday 		string
	Today 			string
}

type Sales struct {
	lastUpdate 		int64
	today 			int
	yesterday 		int
	ratio 			float64
}

type SalesDetailData struct {
	Data    		map[string]int
	lastUpdate 		int64
}
