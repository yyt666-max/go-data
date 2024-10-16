package access

import (
	"fmt"
	"strings"

	"github.com/eolinker/eosc/log"
)

type Access struct {
	Name       string   `yaml:"name" json:"name,omitempty"`
	CName      string   `yaml:"cname" json:"cname,omitempty"`
	Value      string   `yaml:"value" json:"value,omitempty"`
	Apis       []string `yaml:"apis" json:"apis,omitempty"`
	Dependents []string `yaml:"dependents" json:"dependents,omitempty"`
	Children   []Access `yaml:"children" json:"children,omitempty"`
	GuestAllow bool     `yaml:"guest_allow" json:"guest_allow,omitempty"`
}

var (
	access = make(map[string][]Access)
)

func All() map[string][]Access {
	return access
}
func Get(name string) ([]Access, bool) {
	list, has := access[name]
	return list, has
}

func Add(group string, asl []Access) {
	gl := make([]Access, 0, len(asl))
	group = formatGroup(group)
	gp := fmt.Sprint(group, ".")
	for _, a := range asl {
		a.Name = strings.ToLower(a.Name)
		if !strings.HasPrefix(a.Name, gp) {
			a.Name = fmt.Sprint(gp, a.Name)
		}
		gl = append(gl, a)
	}

	access[group] = append(access[group], gl...)
	permits.Set(group, newPermit(group, gl))
}

type Detail struct {
	GuestAllow bool
	Apis       []string
}

func formatAccess(as []Access) (map[string]*Detail, []Template) {
	result := map[string]*Detail{}
	templates := make([]Template, 0, len(as))
	for _, a := range as {
		template := Template{
			Name:       a.Name,
			CName:      a.CName,
			Value:      a.Value,
			Dependents: a.Dependents,
			Children:   []Template{},
		}
		if a.Children != nil {
			childrenResult, childTemplate := formatAccess(a.Children)
			for k, v := range childrenResult {
				result[fmt.Sprintf("%s.%s", a.Value, k)] = v
			}
			template.Children = childTemplate
		} else {
			result[a.Value] = &Detail{
				GuestAllow: a.GuestAllow,
			}
			if a.Apis != nil {
				apis := make([]string, 0, len(a.Apis))
				for _, api := range a.Apis {
					f, err := formatApi(api)
					if err != nil {
						log.Error(err)
						continue
					}
					apis = append(apis, f)
				}
				result[a.Value].Apis = apis
			}
		}

		templates = append(templates, template)
	}
	return result, templates
}

func formatApi(api string) (string, error) {
	index := strings.Index(api, ":")
	if index < 0 {
		return "", fmt.Errorf("api %s format error", api)
	}
	method := strings.TrimSpace(strings.ToUpper(api[:index]))
	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" {
		return "", fmt.Errorf("api %s format error", api)
	}
	path := strings.TrimSpace(api[index+1:])

	return fmt.Sprintf("%s:%s", method, path), nil
}

func formatGroup(group string) string {
	group = strings.ToLower(group)
	group = strings.TrimSpace(group)
	group = strings.Trim(group, ".")
	group = strings.ReplaceAll(group, "-", "_")
	group = strings.ReplaceAll(group, ".", "_")

	return group
}
