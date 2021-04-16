package parser

import (
	"fmt"
	"strings"
)

// mysqlMsgParser is used for parse mysql msg
type mysqlMsgParser struct {}
func MysqlMsgParser() *mysqlMsgParser { return new(mysqlMsgParser) }

// EntryDuplicate method parse & return entry, key from mysql message
func (mp *mysqlMsgParser) EntryDuplicate(msg string) (entry, key string) {
	// Ex) Duplicate entry 'testID' for key 'id' -> Duplicate entry testID for key id
	msg = strings.ReplaceAll(msg, "'", "")

	if _, err := fmt.Sscanf(msg, "Duplicate entry %s for key %s", &entry, &key); err != nil {
		entry, key = "", ""
	}
	return
}

