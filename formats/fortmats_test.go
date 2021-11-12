package formats

import (
	"testing"
)

func Test_getFormats(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name        string
		args        args
		wantFormats []Format
		wantErr     bool
	}{
		{name: "test1", args: args{path: "./formats.json"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFormats, err := loadFormats(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFormats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("got Formats:\n %+v", gotFormats)
		})
	}
}
