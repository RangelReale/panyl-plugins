package parse

import "time"

// parseInLocation parses a timestamp that has no zone information in the given location, or in UTC if nil.
func parseInLocation(layout, value string, loc *time.Location) (time.Time, error) {
	if loc == nil {
		loc = time.UTC
	}
	return time.ParseInLocation(layout, value, loc)
}
