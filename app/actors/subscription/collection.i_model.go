package subscription

import (
	"github.com/ottemo/commerce/app/models"
	"github.com/ottemo/commerce/app/models/subscription"
	"github.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/env"
)

// GetModelName returns model name for the Subscription Collection
func (it *DefaultSubscriptionCollection) GetModelName() string {
	return subscription.ConstModelNameSubscriptionCollection
}

// GetImplementationName returns model implementation name for the Subscription Collection
func (it *DefaultSubscriptionCollection) GetImplementationName() string {
	return "Default" + subscription.ConstModelNameSubscriptionCollection
}

// New returns new instance of model implementation object for the Subscription Collection
func (it *DefaultSubscriptionCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(ConstCollectionNameSubscription)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultSubscriptionCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
