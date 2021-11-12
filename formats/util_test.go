package formats

import "testing"

func Test_subset(t *testing.T) {
	type args struct {
		first  []Tag
		second []Tag
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "a is subset of a,b,c", args: args{first: []Tag{
			"a",
		}, second: []Tag{
			"a",
			"b",
			"c",
		}}, want: true},
		{name: "d is not subset of a,b,c", args: args{first: []Tag{
			"d",
		}, second: []Tag{
			"a",
			"b",
			"c",
		}}, want: false},
		{name: "a,d is not subset of a,b,c", args: args{first: []Tag{
			"a",
			"d",
		}, second: []Tag{
			"a",
			"b",
			"c",
		}}, want: false},
		{name: "a,b is subset of a,b,c", args: args{first: []Tag{
			"a",
			"b",
		}, second: []Tag{
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
