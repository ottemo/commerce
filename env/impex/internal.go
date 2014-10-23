package impex

import (
	"encoding/csv"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/ottemo/foundation/utils"
)

var (
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
	CSV_COLUMN_REGEXP = regexp.MustCompile(`^\s*([~^?])?((?:@?\w+\.)*@?\w+)(\s+(?:=|>)\s*\w+)?(?:\s+<([^>]+)>)?\s*(.*)$`)
)

// converts map[string]interface{} to csv data
func MapToCSV(input []map[string]interface{}, output io.Writer) error {

	// preparing writer for csv file
	//-------------------------------
	csvWriter := csv.NewWriter(output)
	csvWriter.Comma = ','

	columnPath := make(map[string]string)

	var collectColumns func(mapItem map[string]interface{}, path string) = nil
	var getPathValue func(item map[string]interface{}, path []string) interface{}

	collectColumns = func(mapItem map[string]interface{}, path string) {
		for itemKey, itemKeyValue := range mapItem {
			currentPath := strings.Trim(path+"."+itemKey, ".")

			switch typedValue := itemKeyValue.(type) {
			case map[string]interface{}:
				collectColumns(typedValue, currentPath)
			case []map[string]interface{}:
				for _, typedValueListItem := range typedValue {
					collectColumns(typedValueListItem, currentPath)
				}
			case []interface{}:
				isMapItemsInside := false
				for _, typedValueListItem := range typedValue {
					if typedValueListItemAsMap, ok := typedValueListItem.(map[string]interface{}); ok {
						collectColumns(typedValueListItemAsMap, currentPath)
						isMapItemsInside = true
					} else {
						isMapItemsInside = false
						break
					}
				}
				if !isMapItemsInside {
					columnPath[currentPath] = "^" + currentPath
				}
			default:
				columnPath[currentPath] = currentPath
			}
		}
	}

	getPathValue = func(item map[string]interface{}, pathArray []string) interface{} {
		if len(pathArray) == 0 {
			return nil
		}

		if keyValue, present := item[pathArray[0]]; present {

			followingPath := pathArray[1:]
			if keyValueAsList, ok := keyValue.([]interface{}); ok {
				result := make([]interface{}, 0)
				for _, listValue := range keyValueAsList {
					if listValueAsMap, ok := listValue.(map[string]interface{}); ok {
						result = append(result, getPathValue(listValueAsMap, followingPath))
					} else {
						result = append(result, listValue)
					}
				}
				return result
			} else if keyValueAsMap, ok := keyValue.(map[string]interface{}); ok {
				return getPathValue(keyValueAsMap, followingPath)
			}

			if len(pathArray) == 1 {
				return keyValue
			}
		}

		return nil
	}

	// making header
	//---------------
	for _, mapItem := range input { // 1st loop - collecting information for header
		collectColumns(mapItem, "")
	}

	sortedPaths := make([]string, 0, len(columnPath))
	for path, _ := range columnPath {
		sortedPaths = append(sortedPaths, path)
	}
	sort.Strings(sortedPaths)

	csvHeader := make([]string, 0)
	for _, currentPath := range sortedPaths {
		csvHeader = append(csvHeader, columnPath[currentPath])
	}

	csvWriter.Write(csvHeader)
	csvWriter.Flush()

	// making contents
	//------------------
	numberOfColumns := len(columnPath)

	for _, mapItem := range input { // 2nd loop - writing content rows
		// one record by default for item
		itemCSVRecords := make([][]string, 0)
		itemCSVRecords = append(itemCSVRecords, make([]string, numberOfColumns))

		for columnIdx, columnPath := range sortedPaths {
			columnValue := getPathValue(mapItem, strings.Split(columnPath, "."))

			if arrayValue, ok := columnValue.([]interface{}); ok {
				for lineIdx, lineValue := range arrayValue {
					if len(itemCSVRecords) > lineIdx {
						itemCSVRecords[lineIdx][columnIdx] = utils.InterfaceToString(lineValue)
					} else {
						newCSVRecord := make([]string, numberOfColumns)
						newCSVRecord[columnIdx] = utils.InterfaceToString(lineValue)
						itemCSVRecords = append(itemCSVRecords, newCSVRecord)
					}
				}
			} else {
				itemCSVRecords[0][columnIdx] = utils.InterfaceToString(columnValue)
			}
		}

		for _, csvRecord := range itemCSVRecords {
			csvWriter.Write(csvRecord)
		}
		csvWriter.Flush()
	}

	return nil
}

// converts csv data to map[string]interface{} and sends to processorFunc
func CSVToMap(input io.Reader, processorFunc func(item map[string]interface{}) bool) error {

	// preparing reader for csv file
	//-------------------------------
	csvReader := csv.NewReader(input)
	csvReader.Comma = ','

	// reading header/columns information
	//------------------------------------
	csvColumns, err := csvReader.Read()
	if err != nil {
		return err
	}
	csvColumnsNumber := len(csvColumns)

	csvColumnFlags := make([]string, csvColumnsNumber)
	csvColumnPath := make([]string, csvColumnsNumber)
	csvColumnMemorizeAs := make([]string, csvColumnsNumber)
	csvColumnMemorizeTo := make([]string, csvColumnsNumber)
	csvColumnType := make([]string, csvColumnsNumber)
	csvColumnConvertors := make([]string, csvColumnsNumber)

	// extracting column header parts
	for idx, column := range csvColumns {
		regexpGroups := CSV_COLUMN_REGEXP.FindStringSubmatch(column)

		if len(regexpGroups) == 0 { // un-recognized column header
			if strings.TrimSpace(column) != "" { // unless it is blank, considering as path
				csvColumnPath[idx] = column
			}
			continue
		}

		csvColumnFlags[idx] = regexpGroups[1]
		csvColumnPath[idx] = regexpGroups[2]

		tmp := strings.Split(strings.TrimSpace(regexpGroups[3]), " ")
		if len(tmp) == 2 {
			csvColumnMemorizeAs[idx] = strings.TrimSpace(tmp[0])
			csvColumnMemorizeTo[idx] = strings.TrimSpace(tmp[1])
		}

		csvColumnType[idx] = strings.TrimSpace(regexpGroups[4])
		csvColumnConvertors[idx] = strings.TrimSpace(regexpGroups[5])
	}

	// reading CSV contents
	//----------------------
	result := make([]map[string]interface{}, 0)
	csvRecordMap := make(map[string]interface{})
	csvMemorize := make(map[string]interface{})

	// column path can be non static, there we storing currently calculated value
	columnPath := make([]string, csvColumnsNumber)
	columnPathArray := make([][]string, csvColumnsNumber)

	csvRecordNum := 1
	for csvRecord, err := csvReader.Read(); err == nil; csvRecord, err = csvReader.Read() { // csv records loop

		// If path collapse flag is set - we have a new object on that csv row for given path
		// each sub-path have own collapse indicator, if lower path collapses, it not affects top path
		// but not vice versa, if top one collapses, lower should also
		pathCollapseFlag := make(map[string]bool)

		// 1st loop - looking for collapsing paths
		//-----------------------------------------
		blankRecordFlag := true
		for columnIdx, value := range csvRecord { // 1st loop
			// skipping blank values and columns without header information
			if value == "" || csvColumnPath[columnIdx] == "" {
				continue
			}

			blankRecordFlag = false

			isArrayColumn := (csvColumnFlags[columnIdx] == "^")
			isMaybeArrayColumn := (csvColumnFlags[columnIdx] == "?")
			isIgnoreColumn := (csvColumnFlags[columnIdx] == "~")

			keyPath := csvColumnPath[columnIdx]
			keyPathArray := strings.Split(keyPath, ".")

			// updating variable paths (with @) if needed,
			// and storing calculations for following data processing loop
			keyPathWasUpdated := false
			for idx, pathValue := range keyPathArray {
				if strings.HasPrefix(pathValue, "@") {
					if value, present := csvMemorize[pathValue[1:]]; present {
						if stringValue, ok := value.(string); ok {
							keyPathArray[idx] = stringValue
							keyPathWasUpdated = true
						}
					}
				}
			}
			if keyPathWasUpdated {
				keyPath = strings.Join(keyPathArray, ".")
				keyPathArray = strings.Split(keyPath, ".")
			}
			columnPath[columnIdx] = keyPath
			columnPathArray[columnIdx] = keyPathArray

			// storing column value or path+value as variable for variable paths
			switch csvColumnMemorizeAs[columnIdx] {
			case "=":
				csvMemorize[csvColumnMemorizeTo[columnIdx]] = keyPath + "." + value
			case ">":
				csvMemorize[csvColumnMemorizeTo[columnIdx]] = value
			}

			// setting collapse condition for path path if it not array, and value was changed
			if !isArrayColumn && !isMaybeArrayColumn && !isIgnoreColumn && csvColumnMemorizeAs[columnIdx] == "" {
				path := ""
				if arrayLen := len(keyPathArray); arrayLen > 1 {
					path = strings.Join(keyPathArray[:arrayLen-1], ".")
				}
				pathCollapseFlag[path] = true
			}
		}

		// check for blank record
		//------------------------
		if blankRecordFlag {
			break // if we found blank record - this means end data chunk
		}

		// checking for top level collapse
		//---------------------------------
		if value, present := pathCollapseFlag[""]; present && value {
			pathCollapseFlag = make(map[string]bool)

			if csvRecordNum != 1 { // on first row there are blank map
				if !processorFunc(csvRecordMap) {
					return nil // so, processor says to stop
				}
			}
			csvRecordMap = make(map[string]interface{})
		}

		// processing record values (2nd loop)
		//-------------------------------------
		for columnIdx, value := range csvRecord { // 2nd loop
			// skipping blank values, columns without header, path folding columns
			if value == "" || csvColumnPath[columnIdx] == "" || csvColumnMemorizeAs[columnIdx] == "=" {
				continue
			}

			isArrayColumn := (csvColumnFlags[columnIdx] == "^")
			isMaybeArrayColumn := (csvColumnFlags[columnIdx] == "?")

			keyPathArray := columnPathArray[columnIdx]
			lastPathIdx := len(keyPathArray) - 1

			var prevPathMap map[string]interface{} = nil
			var prevPathKey string = ""
			var prevPathValue interface{} = result

			currentPathMap := csvRecordMap
			currentPath := ""

			// moving down through column path
			for idx, key := range keyPathArray { // path loop

				// check we need new object for path
				if value, present := pathCollapseFlag[currentPath]; present && value {
					pathCollapseFlag[currentPath] = false

					newMapValue := make(map[string]interface{})
					if prevPathValueAsList, isList := prevPathValue.([]interface{}); isList {
						prevPathMap[prevPathKey] = append(prevPathValueAsList, newMapValue)
					} else {
						if prevPathValueAsMap, isMap := prevPathValue.(map[string]interface{}); isMap && len(prevPathValueAsMap) > 0 {
							prevPathMap[prevPathKey] = []interface{}{prevPathValue, newMapValue}
						} else {
							prevPathMap[prevPathKey] = newMapValue
						}
					}
					currentPathMap = newMapValue
				}

				if idx == lastPathIdx { // we are at end of key path (i.e. on x for key like a.b.c.d.x)

					currentKeyValue, present := currentPathMap[key]
					if present && (isArrayColumn || isMaybeArrayColumn) {
						if currentValueAsArray, ok := currentKeyValue.([]interface{}); ok {
							currentPathMap[key] = append(currentValueAsArray, value)
						} else {
							currentPathMap[key] = []interface{}{currentKeyValue, value}
						}
					} else {
						if isArrayColumn {
							currentPathMap[key] = []interface{}{value}
						} else {
							currentPathMap[key] = value
						}
					}

				} else { // still looping through path (i.e. on i for path like i.i.i.i.x)

					prevPathMap = currentPathMap
					prevPathKey = key

					currentKeyValue, present := currentPathMap[key]
					prevPathValue = currentKeyValue

					if currentKeyValueAsMap, isMap := currentKeyValue.(map[string]interface{}); present && isMap {
						currentPathMap = currentKeyValueAsMap
					} else if currentKeyValueAsList, isList := currentKeyValue.([]interface{}); present && isList {
						makeNewMapFlag := true
						if len(currentKeyValueAsList) > 0 {
							prevPathValue = currentKeyValueAsList
							lastItemValue := currentKeyValueAsList[len(currentKeyValueAsList)-1]
							if lasItemValueAsMap, isMap := lastItemValue.(map[string]interface{}); isMap {
								currentPathMap = lasItemValueAsMap
								makeNewMapFlag = false
							}
						}
						if makeNewMapFlag {
							newMapValue := make(map[string]interface{})
							prevPathValue = append(currentKeyValueAsList, newMapValue)
							currentPathMap[key] = prevPathValue
							currentPathMap = newMapValue
						}
					} else {
						newMapValue := make(map[string]interface{})
						prevPathValue = newMapValue
						currentPathMap[key] = newMapValue
						currentPathMap = newMapValue
					}
				}

				// updating path
				if currentPath != "" {
					currentPath += "."
				}
				currentPath += key
			}
		}

		csvRecordNum += 1
	}
	processorFunc(csvRecordMap)

	return nil
}
