package rts

import "time"

/**
 * referrers = {
 * 		'url_1': N,
 * 		'url_2': N
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
 * topSellers = {
 *		"lastUpdate": timestamp,
 *		"Data": {
 *			"itemID_1": {
 *				"Name": "XXX",
 *				"Image": "YYY",
 *				"Count": X,
 *			},
 *			...
 *			"itemID_N": {
 *				"Name": "XXX",
 *				"Image": "YYY",
 *				"Count": X,
 *			},
 *		},
 * }
 *
 * type OnlineReferer struct {
 * 		type int
 * 		time int
 * }
 *
 * OnlineSessions map[string]OnlineReferer
 * OnlineDirect int
 * OnlineSite int
 * OnlineSearch int
 *
 *
 */

const (
	COLLECTION_NAME_SALES_HISTORY = "rts_sales_history"
	COLLECTION_NAME_SALES         = "rts_sales"
	COLLECTION_NAME_VISITORS      = "rts_visitors"

	REFERRER_TYPE_DIRECT = 0
	REFERRER_TYPE_SITE   = 1
	REFERRER_TYPE_SEARCH = 2

	VISITOR_ADD_TO_CART = 1
	VISITOR_CHECKOUT    = 2
	VISITOR_SALES       = 3

	VISITOR_ONLINE_SECONDS  = 10

)

var (
	referrers             = make(map[string]int)
	visitorsInfoToday     = new(dbVisitorRow)
	visitorsInfoYesterday = new(dbVisitorRow)
	sales                 = Sales{}
	salesDetail           = make(map[string]*SalesDetailData)
	topSellers            = new(TopSellers)

	OnlineSessions    = make(map[string]*OnlineReferrer)
	OnlineDirect      = 0
	OnlineSite        = 0
	OnlineSearch      = 0
	OnlineSessionsMax = 0
	OnlineDirectMax   = 0
	OnlineSiteMax     = 0
	OnlineSearchMax   = 0

	searchEngines = []string{"www.daum.net", "www.google.com", "www.eniro.se", "www.naver.com", "www.yahoo.com",
		"www.msn.com", "www.bing.com", "www.aol.com", "www.aol.com", "www.lycos.com", "www.ask.com", "www.altavista.com",
		"search.netscape.com", "www.cnn.com", "www.about.com", "www.mamma.com", "www.alltheweb.com", "www.voila.fr",
		"search.virgilio.it", "www.bing.com", "www.baidu.com", "www.alice.com", "www.yandex.com", "www.najdi.org.mk",
		"www.aol.com", "www.mamma.com", "www.seznam.cz", "www.search.com", "www.wp.pl", "online.onetcenter.org",
		"www.szukacz.pl", "www.yam.com", "www.pchome.com", "www.kvasir.no", "sesam.no", "www.ozu.es", "www.terra.com",
		"www.mynet.com", "www.ekolay.net", "www.rambler.ru"}
)

type Visits struct {
	Data      map[string]map[string]int32
	Yesterday string
	Today     string
}

type VisitorDetail struct {
	Time     time.Time
	Checkout int
}

type dbVisitorRow struct {
	Id       string
	Day      time.Time
	Visitors int
	Cart     int
	Checkout int
	Sales    int
	Details  map[string]*VisitorDetail
}

type Sales struct {
	lastUpdate int64
	today      int
	yesterday  int
	ratio      float64
}

type SalesDetailData struct {
	Data       map[string]int
	lastUpdate int64
}

type TopSellers struct {
	Data       map[string]*SellerInfo
	lastUpdate int64
}

type SellerInfo struct {
	Name  string
	Image string
	Count int
}

type OnlineReferrer struct {
	referrerType int
	time         time.Time
}
