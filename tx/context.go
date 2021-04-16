package tx

import "context"

// Context is interface have func about get & set TX value with embedded context.Context
type Context interface {
	context.Context
	Tx() (tx interface{})
	SetTx(tx interface{})
}
