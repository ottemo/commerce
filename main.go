package main

import (
	"fmt"

	"github.com/ottemo/foundation/app"

	_ "github.com/ottemo/foundation/env/config"
	_ "github.com/ottemo/foundation/env/ini"

	//_ "github.com/ottemo/foundation/database/sqlite"
	_ "github.com/ottemo/foundation/database/mongodb"

	_ "github.com/ottemo/foundation/media/fsmedia"

	_ "github.com/ottemo/foundation/api/rest"

	_ "github.com/ottemo/foundation/models/custom_attributes"
	_ "github.com/ottemo/foundation/models/product/default_product"
	_ "github.com/ottemo/foundation/models/category/default_category"
	_ "github.com/ottemo/foundation/models/visitor/default_visitor"
	_ "github.com/ottemo/foundation/models/visitor/default_address"
)

func main() {
	if err := app.Init(); err != nil {
		fmt.Println(err.Error())
	}

	if err := app.Start(); err != nil {
		fmt.Println(err.Error())
	}

	app.Serve()
}
