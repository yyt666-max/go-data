package auto

import (
	"context"
	"encoding/json"
	"testing"
)

type ListCompleteService struct{}

func (l *ListCompleteService) GetLabels(ctx context.Context, ids ...string) map[string]string {
	m := make(map[string]string)
	for _, id := range ids {
		m[id] = "list:" + id
	}
	return m
}

type UserCompleteService struct{}

func (l *UserCompleteService) GetLabels(ctx context.Context, ids ...string) map[string]string {
	m := make(map[string]string)
	for _, id := range ids {
		m[id] = "user:" + id
	}
	return m
}

func TestCreate(t *testing.T) {
	type Test struct {
		UpdaterId Label    `json:"updater_id" aolabel:"user"`
		Creator   Label    `json:"creator" aolabel:"user"`
		Ptr       *Label   `json:"ptr" label:"ptr"`
		NilT      *Label   `json:"nil_t,omitempty" aolabel:"user"`
		List      []*Label `json:"list" label:"list"`
	}
	RegisterService("user", &UserCompleteService{})
	RegisterService("list", &ListCompleteService{})
	tv := &Test{

		Ptr:       UUIDP("1"),
		UpdaterId: UUID("2"),
		Creator:   UUID("3"),
		List: []*Label{
			UUIDP("4"),
			UUIDP("5"),
		},
	}

	CompleteLabels(context.Background(), tv)

	output, err := json.Marshal(tv)
	if err != nil {
		t.Errorf("json.Marshal error:%v", err)
	}
	t.Logf("%s", string(output))

}
