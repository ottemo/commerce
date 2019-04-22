// +build mysql

package basebuild

import (
	// MySQL based database service
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/ottemo/commerce/db/mysql"
)
