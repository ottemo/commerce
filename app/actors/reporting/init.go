package reporting

import (
	"github.com/ottemo/commerce/api"
)

func init() {
	api.RegisterOnRestServiceStart(setupAPI)
}
