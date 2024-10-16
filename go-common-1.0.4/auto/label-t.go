package auto

import (
	"github.com/eolinker/go-common/utils"
)

const (
	unknownOperator = "unknown"
)

type Label struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (o *Label) Set(operators map[string]string) {
	if v, has := operators[o.Id]; has {
		o.Name = v
	} else {
		o.Name = unknownOperator
	}
}

func UUID(id string) Label {
	return Label{Id: id}
}
func UUIDP(id string) *Label {
	if id == "" {
		return nil
	}
	return &Label{Id: id}
}
func List(ids []string) []Label {
	if len(ids) == 0 {
		return nil
	}
	list := make([]Label, 0, len(ids))
	for i := range ids {
		list = append(list, Label{Id: ids[i]})
	}
	return list
}
func ListP(ids []string) []*Label {
	if len(ids) == 0 {
		return nil
	}
	list := make([]*Label, 0, len(ids))
	for i := range ids {
		list = append(list, &Label{Id: ids[i]})
	}
	return list
}

type labelList []*Label

func (is labelList) UUIDS() []string {

	set := utils.NewSet[string]()
	for i := range is {
		set.Set(is[i].Id)
	}

	return set.ToList()
}

func (is labelList) Set(operators map[string]string) {
	if operators == nil {
		operators = make(map[string]string)
	}
	for i := range is {
		is[i].Set(operators)
	}
}
