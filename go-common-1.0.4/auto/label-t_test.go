package auto

import "testing"

func TestAuto(t *testing.T) {
	type TestType struct {
		TagValue string `aovalue:"operator"`
	}
	type args struct {
		tagValue string
		value    string
		target   any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				tagValue: "operator",
				value:    "test Operator",
				target: &TestType{
					TagValue: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Auto(tt.args.tagValue, tt.args.value, tt.args.target)
			t.Logf("%+v", tt.args.target)
		})
	}
}
