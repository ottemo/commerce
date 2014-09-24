package block

import (
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/db"
)

// returns database collection
func (it *DefaultCMSBlockCollection) GetDBCollection() db.I_DBCollection {
	return it.listCollection
}

// returns list of cms block model items
func (it *DefaultCMSBlockCollection) ListCMSBlocks() []cms.I_CMSBlock {
	result := make([]cms.I_CMSBlock, 0)

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, recordData := range dbRecords {
		cmsBlockModel, err := cms.GetCMSBlockModel()
		if err != nil {
			return result
		}
		cmsBlockModel.FromHashMap(recordData)

		result = append(result, cmsBlockModel)
	}

	return result
}
