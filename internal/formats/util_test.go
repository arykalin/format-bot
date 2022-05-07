package formats

import (
	"testing"

	"github.com/arykalin/format-bot/internal/formats/data_getter"
)

func Test_subset(t *testing.T) {
	type args struct {
		first  []data_getter.Tag
		second []data_getter.Tag
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "a is subset of a,b,c", args: args{first: []data_getter.Tag{
			"a",
		}, second: []data_getter.Tag{
			"a",
			"b",
			"c",
		}}, want: true},
		{name: "d is not subset of a,b,c", args: args{first: []data_getter.Tag{
			"d",
		}, second: []data_getter.Tag{
			"a",
			"b",
			"c",
		}}, want: false},
		{name: "a,d is not subset of a,b,c", args: args{first: []data_getter.Tag{
			"a",
			"d",
		}, second: []data_getter.Tag{
			"a",
			"b",
			"c",
		}}, want: false},
		{name: "a,b is subset of a,b,c", args: args{first: []data_getter.Tag{
			"a",
			"b",
		}, second: []data_getter.Tag{
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
