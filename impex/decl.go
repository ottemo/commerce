// Package impex is a implementation of import/export service
package impex

import (
	"regexp"
)

// Package global variables
var (
	ConstImpexLog = true  // flag indicates to make log of values going to be processed
	ConstDebugLog = false // flag indicates to have extra log information

	/*
	 *	column format: [flags]path [memorize] [type] [convertors]
	 *
	 *	flags - optional column modificator
	 *		format: [~|^|?]
	 *		"~" - ignore column on collapse lookup
	 *		"^" - array column
	 *		"?" - maybe array column
	 *
	 *	path - attribute name in result map sub-levels separated by "."
	 *		format: [@a.b.c.]d
	 *		"@a" - memorized value
	 *
	 *	memorize - marks column to hold value in memorize map, these values can be used in path like "item.@value.label"
	 *		format: ={name} | >{name}
	 *		{name}  - alphanumeric value
	 *		={name} - saves {column path} + {column value} to memorize map
	 *		>{name}	- saves {column value} to memorize map
	 *
	 *	type - optional type for column
	 *		format: <{type}>
	 *		{type} - int | float | bool
	 *
	 *	convertors - text template modifications you can apply to value before use it
	 *		format: see (http://golang.org/pkg/text/template/)
	 */
	ConstCSVColumnRegexp = regexp.MustCompile(`^\s*([~^?])?((?:@?\w+\.)*@?\w+)(\s+(?:=|>)\s*\w+)?(?:\s+<([^>]+)>)?\s*(.*)$`)

	// set of service import commands
	importCmd = make(map[string]InterfaceImpexImportCmd)
)
