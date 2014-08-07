package main

import (
	"fmt"

	"github.com/ottemo/foundation/app"

	_ "github.com/ottemo/foundation/env/config"
	_ "github.com/ottemo/foundation/env/ini"

	//_ "github.com/ottemo/foundation/db/sqlite"
	_ "github.com/ottemo/foundation/db/mongo"

	_ "github.com/ottemo/foundation/media/fsmedia"

	_ "github.com/ottemo/foundation/api/rest"

	_ "github.com/ottemo/foundation/app/actors/cart"
	_ "github.com/ottemo/foundation/app/actors/category"
	_ "github.com/ottemo/foundation/app/actors/product"
	_ "github.com/ottemo/foundation/app/actors/visitor"
	_ "github.com/ottemo/foundation/app/actors/visitor/address"

	_ "github.com/ottemo/foundation/app/actors/checkout"
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
