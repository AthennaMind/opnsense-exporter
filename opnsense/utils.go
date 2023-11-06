package opnsense

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// parseStringToInt parses a string value to an int value.
// The endpoint is used to identify the EndpointPath that the caller used.
// so we can propagate in the *APICallError.
func parseStringToInt(value string, endpoint EndpointPath) (int, *APICallError) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, &APICallError{
			Endpoint:   string(endpoint),
			Message:    fmt.Sprintf("error parsing %s to int: %s", value, err.Error()),
			StatusCode: 0,
		}
	}
	return intValue, nil
}

// parseStringToFloatWithReplace parses a string value to a float64 value.
// The replace pattern is used to remove any characters that are not part of the float64 value.
// The regex is first used to check if the value matches the regex format.
func parseStringToFloatWithReplace(value string, regex *regexp.Regexp, replacePattern string, valueTypeName string, logger log.Logger) float64 {
	if regex.MatchString(value) {
		cleanValue := strings.Replace(value, replacePattern, "", -1)
		parsedValue, err := strconv.ParseFloat(cleanValue, 64)
		if err != nil {
			level.Warn(logger).
				Log(
					"msg", fmt.Sprintf("parsing %s: '%s' to float64 failed", valueTypeName, value),
					"err", err,
				)
			return -1.0
		}
		return parsedValue
	}

	level.Warn(logger).
		Log(
			"msg", fmt.Sprintf("parsing %s: '%s' to float64 failed. Pattern matching failed.", valueTypeName, value),
		)
	return -1.0
}
