package default_product


func (it *DefaultProductModel) GetSku() string  { return it.Sku }
func (it *DefaultProductModel) GetName() string { return it.Name }

func (it *DefaultProductModel) GetDescription() string { return it.Description }

func (it *DefaultProductModel) GetDefaultImage() string { return it.DefaultImage }

func (it *DefaultProductModel) GetPrice() float64 { return it.Price }




