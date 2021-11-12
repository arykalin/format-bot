package formats

import (
	"reflect"
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
			f := formats{}
			err := f.loadJson(tt.args.path)
			if err != nil {
				t.Fatal(err)
			}
			err = f.loadData()
			if err != nil {
				t.Fatal(err)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFormats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("got Formats:\n %+v", f.formats)
			t.Logf("got questions:\n %+v", f.questions)
		})
	}
}

func TestNewFormats(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *formats
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFormats(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFormats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFormats() got = %v, want %v", got, tt.want)
			}
		})
	}
}
