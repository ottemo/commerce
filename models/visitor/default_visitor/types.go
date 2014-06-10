package default_visitor

const (
	VISITOR_COLLECTION_NAME = "visitor"
)


type DefaultVisitor struct {
	id string

	Email   string
	Fname   string
	Lname   string

	BillingAddressId  string
	ShippingAddressId string
}
