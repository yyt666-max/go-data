package access

import (
	"fmt"

	"github.com/eolinker/go-common/permit"

	"github.com/eolinker/eosc"
)

var (
	permits = eosc.BuildUntyped[string, *permitAccess]()
)

type Template struct {
	Name       string     `yaml:"name" json:"name,omitempty"`
	CName      string     `yaml:"cname" json:"cname,omitempty"`
	Value      string     `yaml:"value" json:"value,omitempty"`
	Children   []Template `yaml:"children" json:"children,omitempty"`
	Dependents []string   `yaml:"dependents" json:"dependents,omitempty"`
}

type permitAccess struct {
	group string
	// permits 当前权限下的API列表
	permits eosc.Untyped[string, []string]
	// access api对应需要的权限
	access      eosc.Untyped[string, string]
	guestAccess []string
	// template 模版
	template []Template
}

func newPermit(group string, access []Access) *permitAccess {
	p := &permitAccess{
		group:       group,
		permits:     eosc.BuildUntyped[string, []string](),
		access:      eosc.BuildUntyped[string, string](),
		guestAccess: make([]string, 0),
		template:    nil,
	}
	p.Add(access)
	return p
}

func (p *permitAccess) Valid(access string) error {
	_, has := p.access.Get(access)
	if !has {
		return fmt.Errorf("permitAccess %s not found", access)
	}
	return nil
}

func (p *permitAccess) Add(as []Access) error {
	result, templates := formatAccess(as)
	guestAccess := make([]string, 0)
	for k, vs := range result {
		k = fmt.Sprintf("%s.%s", p.group, k)
		apis := vs.Apis
		p.permits.Set(k, apis)
		for _, v := range apis {
			p.access.Set(v, k)
		}
		if vs.GuestAllow {
			guestAccess = append(guestAccess, k)
		}
		permit.AddPermitRule(k, apis...)
	}
	p.template = templates
	p.guestAccess = guestAccess
	return nil
}

func (p *permitAccess) GetTemplate() []Template {
	return p.template
}

func (p *permitAccess) GetPermits(access string) ([]string, error) {
	perms, has := p.permits.Get(access)
	if !has {
		return nil, fmt.Errorf("permitAccess %s not found", access)
	}
	return perms, nil
}

func (p *permitAccess) GuestAccess() []string {
	return p.guestAccess
}

func (p *permitAccess) AccessKeys() []string {
	return p.permits.Keys()
}

func GetPermit(group string) (*permitAccess, bool) {
	return permits.Get(group)
}

func GuestAccess(group string) []string {
	p, has := permits.Get(group)
	if !has {
		return nil
	}
	return p.GuestAccess()
}
