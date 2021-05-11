package validate

import "regexp"

const (
	parentUUIDRegexString = "^p\\d{10}$"
	itemUUIDRegexString   = "^e\\d{10}$"
	childrenRegexString   = "^c\\d{10}$"
)

var (
	parentUUIDRegex = regexp.MustCompile(parentUUIDRegexString)
	itemUUIDRegex   = regexp.MustCompile(itemUUIDRegexString)
	childrenRegex   = regexp.MustCompile(childrenRegexString)
)
