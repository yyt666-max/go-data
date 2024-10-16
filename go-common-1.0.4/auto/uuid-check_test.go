package auto

import "testing"

func TestUUIDCheck(t *testing.T) {
	type Create struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Master     string `json:"master" aocheck:"user"`
		Department string `json:"department" aocheck:"department"`
	}

	createS := []*Create{
		{
			ID:         "1",
			Name:       "test",
			Master:     "1",
			Department: "1",
		},
		{
			ID:         "2",
			Name:       "test",
			Master:     "",
			Department: "2",
		},
	}

	createCheck := searchIDCheck[[]*Create](createS)
	t.Log(createCheck)
}
