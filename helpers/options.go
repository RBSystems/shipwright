package helpers

import (
	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/structs"
)

// GetIcons gets a list of icons from the database
func GetIcons() ([]string, *nerr.E) {
	iconList, err := db.GetDB().GetIcons()
	if err != nil {
		return iconList, nerr.Translate(err).Add("failed to get icon list")
	}

	return iconList, nil
}

// GetTemplates gets a list of templates from the database
func GetTemplates() ([]structs.Template, *nerr.E) {
	templateList, err := db.GetDB().GetAllTemplates()
	if err != nil {
		return templateList, nerr.Translate(err).Add("failed to get all templates")
	}

	return templateList, nil
}