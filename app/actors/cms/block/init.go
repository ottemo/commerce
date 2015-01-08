package block

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// init makes package self-initialization routine
func init() {
	cmsBlockInstance := new(DefaultCMSBlock)
	var _ cms.InterfaceCMSBlock = cmsBlockInstance
	models.RegisterModel(cms.ConstModelNameCMSBlock, cmsBlockInstance)

	cmsBlockCollectionInstance := new(DefaultCMSBlockCollection)
	var _ cms.InterfaceCMSBlockCollection = cmsBlockCollectionInstance
	models.RegisterModel(cms.ConstModelNameCMSBlockCollection, cmsBlockCollectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)

	templateFunc := func(identifier string) string {
		var recursionSaveFunc func(identifier string) string
		var stack []string

		recursionSaveFunc = func(identifier string) string {
			// prevents infinite loop
			for _, stackIdentifier := range stack {
				if stackIdentifier == identifier {
					return ""
				}
			}
			stack = append(stack, identifier)

			// processing new identifier
			block, err := cms.LoadCMSBlockByIdentifier(identifier)
			if err != nil {
				return ""
			}

			blockContents := block.GetContent()
			if strings.Contains("{{", blockContents) {
				// we are replacing template "block" function with own to prevent recursion
				templateFunctions := utils.GetTemplateFunctions()
				templateFunctions["block"] = recursionSaveFunc

				textTemplate, err := template.New("recursiveParsing").Funcs(templateFunctions).Parse(blockContents)
				if err != nil {
					return blockContents
				}

				var result bytes.Buffer
				err = textTemplate.Execute(&result, block.ToHashMap())
				if err != nil {
					return blockContents
				}

				return result.String()
			}

			return blockContents
		}

		return recursionSaveFunc(identifier)
	}
	utils.RegisterTemplateFunction("block", templateFunc)

}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(ConstCmsBlockCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("identifier", "varchar(255)", true)
	collection.AddColumn("content", "text", false)
	collection.AddColumn("created_at", "datetime", false)
	collection.AddColumn("updated_at", "datetime", false)

	return nil
}
