package formats

import "testing"

func Test_subset(t *testing.T) {
	type args struct {
		first  []string
		second []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "a is subset of a,b,c", args: args{first: []string{
			"a",
		}, second: []string{
			"a",
			"b",
			"c",
		}}, want: true},
		{name: "d is not subset of a,b,c", args: args{first: []string{
			"d",
		}, second: []string{
			"a",
			"b",
			"c",
		}}, want: false},
		{name: "a,d is not subset of a,b,c", args: args{first: []string{
			"a",
			"d",
		}, second: []string{
			"a",
			"b",
			"c",
		}}, want: false},
		{name: "a,b is subset of a,b,c", args: args{first: []string{
			"a",
			"b",
		}, second: []string{
			"a",
			"b",
			"c",
		}}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := subset(tt.args.first, tt.args.second)
			if got != tt.want {
				t.Errorf("subset() = %v, want %v", got, tt.want)
			}
		})
	}
}
