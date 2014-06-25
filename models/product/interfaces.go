package product

import("github.com/ottemo/foundation/models")

type IProduct interface {

	  GetSku() string
	 GetName() string

	GetPrice() float64

	models.IModel
	models.IObject
	models.IStorable
	models.IMapable

	models.ICustomAttributes
}
