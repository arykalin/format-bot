package data_getter

import (
	"testing"
)

func Test_sheetData_getFormats(t *testing.T) {
	type fields struct {
		sheetID    string
		secretPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test1", fields: fields{sheetID: "19jaoKkiLRKH9HZ--HX2sE68LGIKt45alnwfrvxwpJNg", secretPath: "./client_secret.json"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newSheetData(tt.fields.sheetID, tt.fields.secretPath)
			gotFormats, err := s.getFormats()
			if (err != nil) != tt.wantErr {
				t.Errorf("getFormats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("got Formats:\n %+v", gotFormats)
		})
	}
}
