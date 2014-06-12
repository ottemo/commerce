package visitor

const (
	VISITOR_COLLECTION_NAME = "visitor"
)

type DefaultVisitor struct {
	id string

	Email     string
	FirstName string
	LastName  string

	BillingAddress  IVisitorAddress
	ShippingAddress IVisitorAddress
}
