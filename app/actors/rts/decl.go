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

  // INIT TABLE_1
 * select _id from products
 	for date = "2014-01-01" to NOW()
 		DELETE from t1 where date=date
	   	for (_id)
			 select qty from order_items where order_id IN (select _id from orders where created_at > "date") and product_id = "_id"
				SUM(qty)

			 insert (SUM,DATE,PRODUCT_ID) in t1
	   	}
    }

 * select _id from orders where created_at = ""
 *  100   2014-10-10 1 +1

 * _id;date;productId;count (rts_sales_history as t1)

 * 		select count from t1 where productId="x"
 *     		SUM(count) -> t2
 * _id;productId;count;range      (rts_sales as t2)
 * 1  ; xxxxxxxx;  12 ; 2014-01-01:
 * 1  ; xxxxxxxx;  12 ; 12/31/88:10/29/14
 */

const (
	COLLECTION_NAME_SALES_HISTORY = "rts_sales_history"
	COLLECTION_NAME_SALES = "rts_sales"
)

var (
	referrers   = make(map[string]*ReferrerData)
	visits      = Visits{Data: make(map[string]map[string]int32)}
	conversions = make(map[string]map[string]int)
	sales       = Sales{}
	salesDetail = make(map[string]*SalesDetailData)
	topSellers  = new (TopSellers)
)

type ReferrerData struct {
	Data  map[string]map[string]bool
	Count int
}

type Visits struct {
	Data      map[string]map[string]int32
	Yesterday string
	Today     string
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
	Image  string
	Count int
}
