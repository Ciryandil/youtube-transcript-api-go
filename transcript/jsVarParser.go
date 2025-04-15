package transcript

import (
	"encoding/json"
	"fmt"
	"strings"
)

func getStringValueAndStartingIndex(variableName string, rawHtml string) ([]rune, uint16, error) {
	stringParts := strings.Split(rawHtml, fmt.Sprintf("var %s", variableName))
	if len(stringParts) < 2 {
		return nil, 0, fmt.Errorf("starting Index Locator: Data unparsable")
	}
	startingIndex := -1
	runesArr := []rune(stringParts[1])
	for itr, char := range runesArr {
		if char == '{' {
			startingIndex = itr
			break
		}
	}
	if startingIndex == -1 {
		return nil, 0, fmt.Errorf("starting Index Locator: Data unparsable: Character not found")
	}
	return runesArr, uint16(startingIndex), nil
}

func findVarSubstring(runesArr []rune, startingIndex uint16) (string, error) {
	escaped := false
	inQuotes := false
	depth := uint8(1)
	length := uint16(1)
	for i := startingIndex + 1; i < uint16(len(runesArr)); i++ {
		char := runesArr[i]
		if escaped {
			escaped = false
		} else if char == '\\' {
			escaped = true
		} else if char == '"' {
			inQuotes = !inQuotes
		} else if !inQuotes {
			if char == '{' {
				depth += 1
			} else if char == '}' {
				depth -= 1
			}
		}
		length += 1
		if depth == 0 {
			break
		}
	}
	if depth != 0 {
		return "", fmt.Errorf("substring extractor: Data unparsable")
	}
	return string(runesArr[startingIndex:(startingIndex + length)]), nil
}

func parseJSVars(variableName string, rawHtml string) (map[string]interface{}, error) {
	runesArr, startingIndex, err := getStringValueAndStartingIndex(variableName, rawHtml)
	if err != nil {
		return nil, fmt.Errorf("js Var Parsing failed : %v", err)
	}
	substr, err := findVarSubstring(runesArr, startingIndex)
	if err != nil {
		return nil, fmt.Errorf("js Var Parsing failed : %v", err)
	}
	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(substr), &jsonMap)
	if err != nil {
		return nil, fmt.Errorf("js Var Parsing failed: failed to load json: %v", err)
	}
	return jsonMap, nil
}
