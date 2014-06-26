package main

import (
	"fmt"

	"github.com/ottemo/foundation/app"

	_ "github.com/ottemo/foundation/config/default_config"
	_ "github.com/ottemo/foundation/config/default_ini_config"

	//_ "github.com/ottemo/foundation/database/sqlite"
	_ "github.com/ottemo/foundation/database/mongodb"

	_ "github.com/ottemo/foundation/media/fsmedia"

	_ "github.com/ottemo/foundation/rest_service"
	_ "github.com/ottemo/foundation/rest_service/default_rest_service"

	_ "github.com/ottemo/foundation/models/custom_attributes"
	_ "github.com/ottemo/foundation/models/product/default_product"
	_ "github.com/ottemo/foundation/models/category/default_category"
	_ "github.com/ottemo/foundation/models/visitor/default_visitor"
	_ "github.com/ottemo/foundation/models/visitor/default_address"
)

func main() {
	if err := app.Start(); err != nil {
		fmt.Println(err.Error())
	}

	app.Serve()
}
