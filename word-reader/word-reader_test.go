package wordreader

import (
	"path"
	"path/filepath"
	"testing"
)

func TestWordReader_readLines(t *testing.T) {
	workingDir, _ := filepath.Abs("..")
	type fields struct {
		path    string
		strings []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			"test on a little file",
			fields{path.Join(workingDir, "resources", "dataset", "w8.txt"), []string{}},
			8,
			false,
		},
		{
			"test on a big file",
			fields{path.Join(workingDir, "resources", "dataset", "w16384.txt"), []string{}},
			16384,
			false,
		},
		{
			"invalid path",
			fields{path.Join(workingDir, "resources", "dataset", "idontexist.txt"), []string{}},
			0,
			true,
		},
		{
			"test on a giant file",
			fields{path.Join(workingDir, "resources", "dataset", "w131072.txt"), []string{}},
			131072,
			false,
		},
		{
			"test on an empty file",
			fields{path.Join(workingDir, "resources", "dataset", "w0.txt"), []string{}},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wr := &WordReader{
				path:    tt.fields.path,
				Strings: tt.fields.strings,
			}

			got, err := wr.ReadLines()
			if (err != nil) != tt.wantErr {
				t.Errorf("WordReader.readLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WordReader.readLines() = %v, want %v", got, tt.want)
			}
		})
	}
}
