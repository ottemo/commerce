// +build !sqlite,!mysql,!mongo

package basebuild

import (
	_ "github.com/ottemo/commerce/db/postgres"
)
