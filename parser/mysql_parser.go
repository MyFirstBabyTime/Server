package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// mysqlMsgParser is used for parse mysql msg
type mysqlMsgParser struct {}
func MysqlMsgParser() *mysqlMsgParser { return new(mysqlMsgParser) }

// EntryDuplicate method parse & return entry, key from mysql entry duplicate error message
func (mp *mysqlMsgParser) EntryDuplicate(msg string) (entry, key string) {
	// Ex) Duplicate entry 'testID' for key 'id' -> Duplicate entry testID for key id
	msg = regexp.MustCompile("'.*?'").ReplaceAllStringFunc(msg, func(s string) string {
		return strings.ReplaceAll(strings.ReplaceAll(s, "'", ""), " ", "")
	})

	if _, err := fmt.Sscanf(msg, "Duplicate entry %s for key %s", &entry, &key); err != nil {
		entry, key = "", ""
	}
	return
}

// NoReferencedRow method parse & return foreign key from mysql no referenced error message
func (mp *mysqlMsgParser) NoReferencedRow(msg string) string {
	var msgFmt = "Cannot add or update a child row: a foreign key constraint fails %s CONSTRAINT %s FOREIGN KEY %s REFERENCES %s"
	for _, i := range []rune{'`', ',', '(', ')'} {
		msg = strings.ReplaceAll(msg, string(i), "")
	}

	bind := struct {
		fkTable, constraint, fk, refTable string
	}{}
	if _, err := fmt.Sscanf(msg, msgFmt, &bind.fkTable, &bind.constraint, &bind.fk, &bind.refTable); err != nil {
		bind.fk = ""
	}
	return bind.fk
}
