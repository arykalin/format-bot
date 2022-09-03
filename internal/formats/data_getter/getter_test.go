package data_getter

import (
	"testing"
)

func Test_getter_GetData(t *testing.T) {
	type fields struct {
		questionsPath  string
		formatsSheetID string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test1", fields: fields{questionsPath: "./questions.json", formatsSheetID: "1SllTPG_dujOctppRwMFoKA6636OB10To4Gc8HhJMbS8"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &getter{
				questionsPath:  tt.fields.questionsPath,
				formatsSheetID: tt.fields.formatsSheetID,
			}
			gotFormats, gotQuestions, err := g.GetData()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("got Formats:\n %+v", gotFormats)
			t.Logf("got Questions:\n %+v", gotQuestions)
		})
	}
}
