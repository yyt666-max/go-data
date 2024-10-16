package access

import "github.com/eolinker/eosc"

var (
	role = eosc.BuildUntyped[string, []Role]()
)

type Role struct {
	Name       string   `yaml:"name" json:"name,omitempty"`
	CName      string   `yaml:"cname" json:"cname,omitempty"`
	Permits    []string `yaml:"permits" json:"permits,omitempty"`
	Supper     bool     `yaml:"supper" json:"supper,omitempty"`
	Default    bool     `yaml:"default" json:"default,omitempty"`
	Dependents []string `yaml:"dependents" json:"dependents,omitempty"`
}

func RoleAdd(roles map[string][]Role) {
	for k, v := range roles {
		role.Set(k, v)
	}
}

func Roles() map[string][]Role {
	return role.All()
}
