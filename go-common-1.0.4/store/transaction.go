package store

import (
	"context"
	"github.com/eolinker/go-common/autowire"
	"gorm.io/gorm"
	"reflect"
)

var _ ITransaction = (*imlTransaction)(nil)

var TxContextKey = struct{}{}

type ITransaction interface {
	Transaction(ctx context.Context, f func(txCtx context.Context) error) error
}
type imlTransaction struct {
	IDB `autowired:""`
}

// Transaction 执行事务
func (b *imlTransaction) Transaction(ctx context.Context, f func(context.Context) error) error {
	if b.IsTxCtx(ctx) {
		return f(ctx)
	}
	return b.DB(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, TxContextKey, tx)
		return f(txCtx)
	})
}
func init() {
	autowire.Auto[ITransaction](func() reflect.Value {
		return reflect.ValueOf(new(imlTransaction))
	})

}
