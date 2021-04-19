package validate

import "regexp"

const (
	parentUUIDRegexString = "^p\\d{10}$"
)

var (
	parentUUIDRegex = regexp.MustCompile(parentUUIDRegexString)
)
