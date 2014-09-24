package page

import (
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/db"
)

// returns database collection
func (it *DefaultCMSPageCollection) GetDBCollection() db.I_DBCollection {
	return it.listCollection
}

// returns list of cms page model items
func (it *DefaultCMSPageCollection) ListCMSPages() []cms.I_CMSPage {
	result := make([]cms.I_CMSPage, 0)

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, recordData := range dbRecords {
		cmsPageModel, err := cms.GetCMSPageModel()
		if err != nil {
			return result
		}
		cmsPageModel.FromHashMap(recordData)

		result = append(result, cmsPageModel)
	}

	return result
}
