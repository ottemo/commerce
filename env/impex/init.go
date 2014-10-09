package impex

import (
	"github.com/ottemo/foundation/api"
)

func init() {
	api.RegisterOnRestServiceStart(setupAPI)
}
