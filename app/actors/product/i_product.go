package product

func (it *DefaultProduct) GetSku() string  { return it.Sku }
func (it *DefaultProduct) GetName() string { return it.Name }

func (it *DefaultProduct) GetShortDescription() string { return it.ShortDescription }
func (it *DefaultProduct) GetDescription() string { return it.Description }

func (it *DefaultProduct) GetDefaultImage() string { return it.DefaultImage }

func (it *DefaultProduct) GetPrice() float64 { return it.Price }

func (it *DefaultProduct) GetWeight() float64 { return it.Weight }
func (it *DefaultProduct) GetSize() float64 { return it.Size }
