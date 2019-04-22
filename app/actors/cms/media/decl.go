// Package media is a default implementation of cms page related interfaces declared in
// "github.com/ottemo/commerce/app/models/cms" package
package media

import (
	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/media"
)

// Package global constants
const (
	ConstErrorModule = "cms/media"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstStorageModel  = "cms"
	ConstStorageObject = "media"
)

var (
	mediaStorage media.InterfaceMediaStorage
)
