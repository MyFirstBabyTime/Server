package validate

import "regexp"

const (
	parentUUIDRegexString = "^p\\d{10}$"
	itemUUIDRegexString = "^e\\d{10}$"
)

var (
	parentUUIDRegex = regexp.MustCompile(parentUUIDRegexString)
	itemUUIDRegx = regexp.MustCompile(itemUUIDRegexString)
)
