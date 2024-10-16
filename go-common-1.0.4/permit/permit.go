package permit

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/permit/internal/permit_store"
	"github.com/eolinker/go-common/utils"
)

var (
	_ IPermit = (*permitImpl)(nil)
)

func init() {
	autowire.Auto[IPermit](func() reflect.Value {
		return reflect.ValueOf(new(permitImpl))
	})
}

type IPermit interface {
	Check(ctx context.Context, domain string, target []string, access []string) (bool, error)
	Add(ctx context.Context, access string, domain string, target string) error
	Reset(ctx context.Context, access string, domain string, target ...string) error
	Remove(ctx context.Context, access string, domain string, target string) error
	// Delete is a Go function that takes a context, domain string, and variable number of access strings as parameters.
	// It returns an error.
	Delete(ctx context.Context, domain string, access ...string) error
	// GrantSystem is a Go function that takes a context, access, and domain as parameters.
	// It returns a slice of strings and an error.
	Granted(ctx context.Context, access, domain string) ([]string, error)
	// GrantForDomain description of the Go function.
	//
	// ctx context.Context, domain string parameters.
	// Returns map[string][]string, error.
	// map[access][]target
	GrantForDomain(ctx context.Context, domain string) (map[string][]string, error)
	Access(ctx context.Context, domain string, target ...string) ([]string, error)
}

type permitImpl struct {
	store                permit_store.IPermitStore `autowired:""`
	permitInitializeData IPermitInitialize         `autowired:""`
}

func (p *permitImpl) OnComplete() {

	list, err := p.store.List(context.Background(), map[string]interface{}{})
	if err != nil {
		return
	}

	grants := utils.MapChange(utils.SliceToMapArray(list, func(t *permit_store.Permit) string {
		return t.Domain
	}), func(ps []*permit_store.Permit) map[string][]string {
		return utils.SliceToMapArrayO(ps, func(t *permit_store.Permit) (string, string) {
			return t.Access, t.Target
		})
	})

	template := p.permitInitializeData.Grants()

	newGrant := make([]*permit_store.Permit, 0, len(list))
	timeNow := time.Now()
	for domain, tmpM := range template {
		gM, has := grants[domain]
		if !has {
			gM = make(map[string][]string)
		}

		for access, targets := range tmpM {
			if _, has := gM[access]; has {
				continue
			}
			newGrant = append(newGrant, utils.SliceToSlice(targets, func(s string) *permit_store.Permit {
				return &permit_store.Permit{
					Domain:     domain,
					Access:     access,
					Target:     s,
					CreateTime: timeNow,
				}
			})...)
		}

	}
	if len(newGrant) > 0 {
		err := p.store.Insert(context.Background(), newGrant...)
		if err != nil {
			panic(fmt.Sprint("initialize permit :", err))
		}
	}
}

func (p *permitImpl) Check(ctx context.Context, domain string, targets []string, access []string) (bool, error) {
	if len(access) == 0 {
		return true, nil
	}
	if len(targets) == 0 {
		return false, nil
	}
	query, err := p.store.CountQuery(ctx, "`domain` = ? AND `access` in ? AND `target` in ?", domain, access, targets)
	if err != nil {
		return false, err
	}
	return query > 0, nil
}

func (p *permitImpl) Add(ctx context.Context, access string, domain string, target string) error {

	return p.store.Transaction(ctx, func(ctx context.Context) error {

		return p.store.Save(ctx, &permit_store.Permit{
			Id:         0,
			Domain:     domain,
			Access:     access,
			Target:     target,
			CreateTime: time.Now(),
		})

	})
}

func (p *permitImpl) Reset(ctx context.Context, access string, domain string, target ...string) error {
	evl := make([]*permit_store.Permit, 0, len(target))

	nt := time.Now()
	for _, t := range target {
		evl = append(evl, &permit_store.Permit{
			Id:         0,
			Domain:     domain,
			Access:     access,
			Target:     t,
			CreateTime: nt,
		})
	}
	return p.store.Transaction(ctx, func(ctx context.Context) error {
		_, err := p.store.DeleteWhere(ctx, map[string]interface{}{"domain": domain, "access": access})
		if err != nil {
			return err
		}
		for _, v := range evl {
			if err := p.store.Save(ctx, v); err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *permitImpl) Remove(ctx context.Context, access string, domain string, target string) error {
	//_, err := p.store.DeleteQuery(ctx, "`domain` = ? AND `access` = ? AND `target` = ?", domain, access, target)
	_, err := p.store.DeleteWhere(ctx, map[string]interface{}{"domain": domain, "access": access, "target": target})
	return err
}

func (p *permitImpl) Delete(ctx context.Context, domain string, access ...string) error {
	if domain == "" {
		return nil
	}
	if len(access) == 0 {
		//_, err := p.store.DeleteQuery(ctx, "`domain` = ?", domain)
		_, err := p.store.DeleteWhere(ctx, map[string]interface{}{"domain": domain})
		return err
	}
	if len(access) == 1 {
		//_, err := p.store.DeleteQuery(ctx, "`domain` = ? AND `access` = ?", domain, access[0])
		_, err := p.store.DeleteWhere(ctx, map[string]interface{}{"domain": domain, "access": access[0]})
		return err
	}

	//_, err := p.store.DeleteQuery(ctx, "`domain` = ? AND `access` in (?)", domain, access)
	_, err := p.store.DeleteWhere(ctx, map[string]interface{}{"domain": domain, "access": access})
	return err
}

func (p *permitImpl) Granted(ctx context.Context, access string, domain string) ([]string, error) {

	list, err := p.store.ListQuery(ctx, "`domain` = ? AND `access` = ?", []interface{}{domain, access}, "target ASC")
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(list, func(i *permit_store.Permit) string { return i.Target }), nil
}

func (p *permitImpl) GrantForDomain(ctx context.Context, domain string) (map[string][]string, error) {
	list, err := p.store.ListQuery(ctx, "`domain` = ? ", []interface{}{domain}, "target ASC")
	if err != nil {
		return nil, err
	}
	return utils.SliceToMapArrayO(list, func(i *permit_store.Permit) (string, string) { return i.Access, i.Target }), nil
}

func (p *permitImpl) Access(ctx context.Context, domain string, target ...string) ([]string, error) {

	list, err := p.store.List(ctx,
		map[string]interface{}{
			"domain": domain,
			"target": target,
		})

	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(list, func(i *permit_store.Permit) string { return i.Access }), nil

}
func init() {
	autowire.Auto[IPermit](func() reflect.Value {
		return reflect.ValueOf(new(permitImpl))
	})
}

var (
	_ IPermit = (*permitImpl)(nil)
)
