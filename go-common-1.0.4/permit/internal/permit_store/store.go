package permit_store

import (
	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/store"
	"reflect"
)

type IPermitStore interface {
	store.IBaseStore[Permit]
}

type imlPermitStore struct {
	store.Store[Permit]
}

func init() {
	autowire.Auto[IPermitStore](func() reflect.Value {
		return reflect.ValueOf(new(imlPermitStore))
	})
}
