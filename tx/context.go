package tx

import "context"

// Context is interface have func about get & set TX value with embedded context.Context
type Context interface {
	context.Context
	Tx() (tx interface{})
	SetTx(tx interface{})
}

// txContext is implementation of Context interface
type txContext struct {
	context.Context
	txKey interface{}
}

// Tx method Get TX value from context
func (tc *txContext) Tx() (tx interface{}) { return tc.Value(tc.txKey) }

// SetTx method Set TX value in context
func (tc *txContext) SetTx(tx interface{}) { tc.Context = context.WithValue(tc.Context, tc.txKey, tx) }
