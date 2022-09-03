package formats

import (
	"testing"

	"github.com/arykalin/format-bot/internal/formats/data_getter"
)

func TestNewFormats(t *testing.T) {
	type args struct {
		path    string
		sheetID string
	}
	tests := []struct {
		name    string
		args    args
		want    *formats
		wantErr bool
	}{
		{name: "test1", args: args{path: "./data_getter/questions.json", sheetID: "xxx"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getter := data_getter.NewGetter(tt.args.path, tt.args.sheetID)
			got, err := NewFormats(getter)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFormats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			f, err := got.GetFormats(nil)
			t.Logf("got formats %+v", f)
		})
	}
}
