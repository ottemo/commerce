package rts

import (
	"time"

	"github.com/ottemo/commerce/api"
	"github.com/ottemo/commerce/app/models/cart"
	"github.com/ottemo/commerce/app/models/order"
	"github.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/utils"
)

func visitsHandler(event env.InterfaceEvent) error {

	if sessionInstance, ok := event.Get("session").(api.InterfaceSession); ok && sessionInstance != nil {
		if sessionID := sessionInstance.GetID(); sessionID != "" {
			currentHour := time.Now().Truncate(time.Hour).Unix()
			CheckHourUpdateForStatistic()

			// Total page views
			statistic[currentHour].TotalVisits++
			monthStatistic.TotalVisits++

			// TODO: Super flakey implementation for telling if the visitor has been tracked today
			// by reusing an 'add to bag' tracking mechanism
			// commerce/app/actors/rts/decl.go :45
			if _, present := visitState[sessionID]; !present {
				visitState[sessionID] = false

				// Unique page views
				statistic[currentHour].Visit++
				monthStatistic.Visit++

				err := SaveStatisticsData()
				if err != nil {
					_ = env.ErrorDispatch(err)
				}
			}
		}
	}

	return nil
}

func addToCartHandler(event env.InterfaceEvent) error {

	if sessionInstance, ok := event.Get("session").(api.InterfaceSession); ok && sessionInstance != nil {
		if sessionID := sessionInstance.GetID(); sessionID != "" {

			currentHour := time.Now().Truncate(time.Hour).Unix()
			CheckHourUpdateForStatistic()

			// Add cart counter if it's a visitor that work in this hour
			if haveCard, present := visitState[sessionID]; present {
				if !haveCard {
					visitState[sessionID] = true

					if _, present := statistic[currentHour]; present && statistic[currentHour] != nil {
						statistic[currentHour].Cart++
						monthStatistic.Cart++
					}

					err := SaveStatisticsData()
					if err != nil {
						_ = env.ErrorDispatch(err)
					}
				}

				// Add cart and visit counter if it's a visitor that work for a past hour
			} else {
				visitState[sessionID] = true

				if _, present := statistic[currentHour]; present && statistic[currentHour] != nil {
					statistic[currentHour].Visit++
					statistic[currentHour].TotalVisits++
					statistic[currentHour].Cart++
				}

				err := SaveStatisticsData()
				if err != nil {
					_ = env.ErrorDispatch(err)
				}
			}
		}
	}

	return nil
}

func visitCheckoutHandler(event env.InterfaceEvent) error {

	currentHour := time.Now().Truncate(time.Hour).Unix()
	CheckHourUpdateForStatistic()

	if _, present := statistic[currentHour]; present && statistic[currentHour] != nil {
		statistic[currentHour].VisitCheckout++
		monthStatistic.VisitCheckout++
	}

	err := SaveStatisticsData()
	if err != nil {
		_ = env.ErrorDispatch(err)
	}

	return nil
}

func setPaymentHandler(event env.InterfaceEvent) error {

	currentHour := time.Now().Truncate(time.Hour).Unix()
	CheckHourUpdateForStatistic()

	if _, present := statistic[currentHour]; present && statistic[currentHour] != nil {
		statistic[currentHour].SetPayment++
		monthStatistic.SetPayment++
	}

	err := SaveStatisticsData()
	if err != nil {
		_ = env.ErrorDispatch(err)
	}

	return nil
}

func purchasedHandler(event env.InterfaceEvent) error {

	if sessionInstance, ok := event.Get("session").(api.InterfaceSession); ok && sessionInstance != nil {
		if sessionID := sessionInstance.GetID(); sessionID != "" {
			saleAmount := float64(0)
			if orderInstance, ok := event.Get("order").(order.InterfaceOrder); ok {
				saleAmount = orderInstance.GetGrandTotal()
			}

			currentHour := time.Now().Truncate(time.Hour).Unix()
			CheckHourUpdateForStatistic()

			if _, present := visitState[sessionID]; !present {
				// Increasing sales, cart and visit counters for visitor of a purchase
				if _, present := statistic[currentHour]; present && statistic[currentHour] != nil {
					statistic[currentHour].Visit++
					statistic[currentHour].TotalVisits++
					statistic[currentHour].Cart++
					statistic[currentHour].VisitCheckout++
					statistic[currentHour].SetPayment++
					monthStatistic.Visit++
					monthStatistic.TotalVisits++
					monthStatistic.Cart++
					monthStatistic.VisitCheckout++
					monthStatistic.SetPayment++
				}
			}

			visitState[sessionID] = false
			if _, present := statistic[currentHour]; present && statistic[currentHour] != nil {
				statistic[currentHour].Sales++
				statistic[currentHour].SalesAmount += saleAmount
				monthStatistic.Sales++
				monthStatistic.SalesAmount += saleAmount
			}

			err := SaveStatisticsData()
			if err != nil {
				_ = env.ErrorDispatch(err)
			}
		}
	} else {
		if orderInstance, ok := event.Get("order").(order.InterfaceOrder); ok {
			saleAmount := orderInstance.GetGrandTotal()

			currentHour := time.Now().Truncate(time.Hour).Unix()
			CheckHourUpdateForStatistic()

			if _, present := statistic[currentHour]; present && statistic[currentHour] != nil {
				statistic[currentHour].Sales++
				statistic[currentHour].SalesAmount += saleAmount
				monthStatistic.Sales++
				monthStatistic.SalesAmount += saleAmount
			}

			err := SaveStatisticsData()
			if err != nil {
				_ = env.ErrorDispatch(err)
			}
		}
	}

	return nil
}

func salesHandler(event env.InterfaceEvent) error {

	if cartInstance, ok := event.Get("cart").(cart.InterfaceCart); ok && cartInstance != nil {
		cartProducts := cartInstance.GetItems()

		if len(cartProducts) == 0 {
			return nil
		}

		productQty := make(map[string]int)
		for _, cartItem := range cartProducts {
			productQty[cartItem.GetProductID()] += cartItem.GetQty()
		}

		salesHistoryCollection, err := db.GetCollection(ConstCollectionNameRTSSalesHistory)
		if err != nil {
			return env.ErrorDispatch(err)
		}
		currentDate := time.Now().Truncate(time.Hour).Add(time.Hour)
		for productID, count := range productQty {

			salesHistoryRecord := make(map[string]interface{})

			if err := salesHistoryCollection.ClearFilters(); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "121c0d25-de60-47d7-99ee-e393af75825e", err.Error())
			}
			if err := salesHistoryCollection.AddFilter("created_at", "=", currentDate); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "051a6574-108b-41b3-983e-c05da5ba887e", err.Error())
			}
			if err := salesHistoryCollection.AddFilter("product_id", "=", productID); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6405bd35-8410-41ca-9e1f-75eb36b32531", err.Error())
			}
			dbSaleRow, err := salesHistoryCollection.Load()
			if err != nil {
				return env.ErrorDispatch(err)
			}

			//	rewrite existing record if we have one in database
			newCount := utils.InterfaceToInt(count)
			if len(dbSaleRow) > 0 {
				salesHistoryRecord["_id"] = dbSaleRow[0]["_id"]
				newCount = newCount + utils.InterfaceToInt(dbSaleRow[0]["count"])
			}

			// saving new history record
			if newCount > 0 {
				salesHistoryRecord["product_id"] = productID
				salesHistoryRecord["created_at"] = currentDate
				salesHistoryRecord["count"] = newCount
				_, err = salesHistoryCollection.Save(salesHistoryRecord)
				if err != nil {
					return env.ErrorDispatch(err)
				}
			}
		}
	}

	return nil
}
