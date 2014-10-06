package review

import (
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func setupAPI() error {

	var err error = nil

	err = api.GetRestService().RegisterAPI("product", "GET", "review/list/:pid", restReviewList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "review/add/:pid", restReviewAdd)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "review/add/:pid/:stars", restReviewAdd)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "DELETE", "review/remove/:reviewId", restReviewRemove)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "GET", "rating/info/:pid", restRatingInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API function used to get list of reviews for particular product
func restReviewList(params *api.T_APIHandlerParams) (interface{}, error) {

	productObject, err := product.LoadProductById(params.RequestURLParams["pid"])
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection, err := db.GetCollection(REVIEW_COLLECTION_NAME)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection.AddFilter("product_id", "=", productObject.GetId())
	collection.AddFilter("review", "!=", "")

	records, err := collection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return records, nil
}

// WEB REST API function used to add new review for a product
func restReviewAdd(params *api.T_APIHandlerParams) (interface{}, error) {

	visitorObject, err := visitor.GetCurrentVisitor(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	productObject, err := product.LoadProductById(params.RequestURLParams["pid"])
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	reviewCollection, err := db.GetCollection(REVIEW_COLLECTION_NAME)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// rating update if was set
	//-------------------------
	ratingValue := 0
	if starsValue, present := params.RequestURLParams["stars"]; present {

		starsNum := utils.InterfaceToInt(starsValue)
		if starsNum <= 0 || starsNum > 5 {
			return nil, env.ErrorNew("stars should be value integer beetween 1 and 5")
		}

		reviewCollection.AddFilter("product_id", "=", productObject.GetId())
		reviewCollection.AddFilter("visitor_id", "=", visitorObject.GetId())
		reviewCollection.AddFilter("rating", ">", 0)

		records, err := reviewCollection.Count()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		if records > 0 {
			return nil, env.ErrorNew("you have already vote for that product")
		}

		ratingValue = starsNum

		ratingCollection, err := db.GetCollection(RATING_COLLECTION_NAME)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		ratingCollection.AddFilter("product_id", "=", productObject.GetId())
		ratingRecords, err := ratingCollection.Load()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		recordAttribute := "stars_" + utils.InterfaceToString(ratingValue)
		var ratingRecord map[string]interface{} = nil

		if len(ratingRecords) > 0 {
			ratingRecord = ratingRecords[0]

			ratingRecord[recordAttribute] = utils.InterfaceToInt(ratingRecord[recordAttribute]) + 1
		} else {
			ratingRecord = map[string]interface{}{
				"product_id": productObject.GetId(),
				"stars_1":    0,
				"stars_2":    0,
				"stars_3":    0,
				"stars_4":    0,
				"stars_5":    0,
			}

			ratingRecord[recordAttribute] = 1
		}
		ratingCollection.Save(ratingRecord)
	}

	// review add new record
	//----------------------
	storingValues := map[string]interface{}{
		"product_id": productObject.GetId(),
		"visitor_id": visitorObject.GetId(),
		"username":   visitorObject.GetFullName(),
		"rating":     ratingValue,
		"review":     params.RequestContent,
		"created_at": time.Now(),
	}

	newId, err := reviewCollection.Save(storingValues)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	storingValues["_id"] = newId

	return storingValues, nil
}

// WEB REST API function used to remove review for a product
func restReviewRemove(params *api.T_APIHandlerParams) (interface{}, error) {

	reviewId := params.RequestURLParams["reviewId"]

	visitorObject, err := visitor.GetCurrentVisitor(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection, err := db.GetCollection(REVIEW_COLLECTION_NAME)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	reviewRecord, err := collection.LoadById(reviewId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if visitorId, present := reviewRecord["visitor_id"]; present {

		// check rights
		if err := api.ValidateAdminRights(params); err != nil {
			if visitorId != visitorObject.GetId() {
				return nil, env.ErrorDispatch(err)
			}
		}

		// rating update if was set
		//-------------------------
		reviewRating := utils.InterfaceToInt(reviewRecord["rating"])

		if reviewRating > 0 {
			ratingCollection, err := db.GetCollection(RATING_COLLECTION_NAME)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			ratingCollection.AddFilter("product_id", "=", reviewRecord["product_id"])
			ratingRecords, err := ratingCollection.Load()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			var ratingRecord map[string]interface{} = nil

			if len(ratingRecords) > 0 {
				ratingRecord = ratingRecords[0]

				recordAttribute := "stars_" + utils.InterfaceToString(reviewRating)
				ratingRecord[recordAttribute] = utils.InterfaceToInt(ratingRecord[recordAttribute]) - 1
				ratingCollection.Save(ratingRecord)
			}
		}

		// review remove
		//--------------
		collection.DeleteById(reviewId)
	}

	return "ok", nil
}

// WEB REST API function used to get product rating info
func restRatingInfo(params *api.T_APIHandlerParams) (interface{}, error) {

	productObject, err := product.LoadProductById(params.RequestURLParams["pid"])
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	ratingCollection, err := db.GetCollection(RATING_COLLECTION_NAME)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	ratingCollection.AddFilter("product_id", "=", productObject.GetId())
	ratingRecords, err := ratingCollection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return ratingRecords, nil
}
