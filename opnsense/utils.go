package opnsense

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
)

// parseStringToInt parses a string value to an int value.
// The endpoint is used to identify the EndpointPath that the caller used.
// so we can propagate in the *APICallError.
func parseStringToInt(value string, endpoint EndpointPath) (int, *APICallError) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, &APICallError{
			Endpoint:   string(endpoint),
			Message:    fmt.Sprintf("error parsing '%s' to int: %s", value, err.Error()),
			StatusCode: 0,
		}
	}
	return intValue, nil
}

// parseStringToFloatWithReplace parses a string value to a float64 value.
// The replace pattern is used to remove any characters that are not part of the float64 value.
// The regex is first used to check if the value matches the regex format.
func parseStringToFloatWithReplace(value string, regex *regexp.Regexp, replacePattern string, valueTypeName string, logger *slog.Logger) float64 {
	if regex.MatchString(value) {
		cleanValue := strings.ReplaceAll(value, replacePattern, "")
		parsedValue, err := strconv.ParseFloat(cleanValue, 64)
		if err != nil {
			logger.Warn(
				fmt.Sprintf("parsing %s: '%s' to float64 failed", valueTypeName, value),
				"details", err,
			)
			return -1.0
		}
		return parsedValue
	}

	logger.Warn(
		fmt.Sprintf("parsing %s: '%s' to float64 failed. Pattern matching failed.", valueTypeName, value),
	)
	return -1.0
}

// sliceIntToMapStringInt is a helper function to convert a slice of strings to a map of string:int
// The key of the map is the string value in the slice and
// the value of the map is the int value of the string.
// The endpoint is used to identify the EndpointPath that the caller used.
// so we can propagate in the *APICallError.
// Fails if any of the string values in the slice cannot be parsed to an int.
func sliceIntToMapStringInt(strings []string, url EndpointPath) (map[string]int, *APICallError) {
	ints := make(map[string]int)

	for _, str := range strings {
		value, err := parseStringToInt(str, url)
		if err != nil {
			return nil, err
		}
		ints[str] = value
	}

	return ints, nil
}

// parseStringToBool parses a string value to a bool value.
// The value is considered true if it is not "0".
func parseStringToBool(value string) bool {
	return value != "0"
}
