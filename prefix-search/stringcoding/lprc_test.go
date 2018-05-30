package stringcoding

import (
	"testing"
)

func TestLPRC_Retrieval(t *testing.T) {
	type fields struct {
		Epsilon float64
		strings []string
	}
	type args struct {
		u uint64
		l uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"First uncompressed string",
			fields{
				70,
				[]string{"ciao", "cic", "cuz"},
			},
			args{
				uint64(0),
				uint64(16),
			},
			"ci",
			false,
		},
		{
			"First uncompressed string",
			fields{
				70,
				[]string{"ciao", "cic", "cuz"},
			},
			args{
				uint64(2),
				uint64(16),
			},
			"cu",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lprc := NewLPRC(tt.fields.strings, tt.fields.Epsilon)
			for i, s := range tt.fields.strings {
				lprc.add(s, uint64(i))
			}

			got, err := lprc.Retrieval(tt.args.u, tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("LPRC.Retrieval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LPRC.Retrieval() = %v, want %v", got, tt.want)
			}
		})
	}
}
