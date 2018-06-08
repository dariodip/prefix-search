package stringcoding

import (
	"testing"

	"reflect"
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
			"0) error test",
			fields{
				1,
				[]string{"caso", "cat", "cena", "delfino"}, //uccc
			},
			args{
				uint64(3),
				uint64(64),
			},
			"delfino",
			false,
		},
		{
			"1) First uncompressed string",
			fields{
				70,
				[]string{"ciao", "cic", "cuz"}, //ucc
			},
			args{
				uint64(0),
				uint64(16),
			},
			"ci",
			false,
		},
		{
			"2) Third compressed string",
			fields{
				70,
				[]string{"caso", "cic", "cuz"}, //ucc
			},
			args{
				uint64(2),
				uint64(16),
			},
			"cu",
			false,
		},
		{
			"3) Different prefix for all",
			fields{
				70,
				[]string{"asso", "basso", "dasso"}, //ucc
			},
			args{
				uint64(2),
				uint64(16),
			},
			"da",
			false,
		},
		{
			"4) Different prefix for all",
			fields{
				1,
				[]string{"asso", "basso", "dasso"}, //ucc
			},
			args{
				uint64(0),
				uint64(16),
			},
			"as",
			false,
		},
		{
			"5) Second compressed string followed by uncompressed",
			fields{
				1,
				[]string{"casotto", "cisonostatierrori", "cuz"}, //ucu
			},
			args{
				uint64(1),
				uint64(16),
			},
			"ci",
			false,
		},
		{
			"6) Third uncompressed string (different prefixes)",
			fields{
				1,
				[]string{"casotto", "visonostatierrori", "zuz"}, //ucu
			},
			args{
				uint64(2),
				uint64(16),
			},
			"zu",
			false,
		},
		{
			"7) Third uncompressed string",
			fields{
				1,
				[]string{"casotto", "cisonostatierrori", "cuz"}, //ucu
			},
			args{
				uint64(2),
				uint64(24),
			},
			"cuz",
			false,
		},
		{
			"8) Forth compressed string, uncompressed not first",
			fields{
				10,
				[]string{"casotto", "cisonostatimoltierrori", "codice", "console", "cuz"}, //ucucc
			},
			args{
				uint64(3),
				uint64(16),
			},
			"co",
			false,
		},
		{
			"9) Index out of bound",
			fields{
				1,
				[]string{"casotto", "cisonostatierrori", "cuz"},
			},
			args{
				uint64(10),
				uint64(16),
			},
			"",
			true,
		},
		{
			"10) l greater than |Si|",
			fields{
				1,
				[]string{"casotto", "cisonostatierrori", "cuz"},
			},
			args{
				uint64(2),
				uint64(128),
			},
			"cuz",
			false,
		},
		{
			"11) last string less 1 char",
			fields{
				1,
				[]string{"casotto", "cisonostatierrori", "cuz", "delfino"},
			},
			args{
				uint64(3),
				uint64(48),
			},
			"delfin",
			false,
		},
		{
			"12) last string full",
			fields{
				1,
				[]string{"casotto", "cisonostatierrori", "cuz", "delfino"},
			},
			args{
				uint64(3),
				uint64(56),
			},
			"delfino",
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
				t.Errorf("len(got)= %d", len(got))
				t.Errorf("LPRC.Retrieval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLPRC_getStringLength(t *testing.T) {
	type fields struct {
		Epsilon float64
		strings []string
	}
	type args struct {
		i uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint64
		wantErr bool
	}{
		{
			"Uncompressed string",
			fields{
				1.0,
				[]string{"caso", "cena", "delfino"},
			},
			args{
				uint64(0),
			},
			uint64(40),
			false,
		},
		{
			"Compressed string",
			fields{
				1.0,
				[]string{"caso", "cena", "delfino"},
			},
			args{
				uint64(1),
			},
			uint64(40),
			false,
		},
		{
			"Compressed string next to a compressed string",
			fields{
				1.0,
				[]string{"caso", "cena", "delfino"},
			},
			args{
				uint64(2),
			},
			uint64(64),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lprc := NewLPRC(tt.fields.strings, tt.fields.Epsilon)
			for i, s := range tt.fields.strings {
				lprc.add(s, uint64(i))
			}
			got, err := lprc.getStringLength(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("LPRC.getStringLength() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LPRC.getStringLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLPRC_FullPrefixSearch(t *testing.T) {
	type fields struct {
		Epsilon float64
		strings []string
	}
	type args struct {
		prefix string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			"Last string",
			fields{
				1,
				[]string{"casotto", "cisonostatierrori", "cuz", "delfino"},
			},
			args{
				"de",
			},
			[]string{"delfino"},
			false,
		},
		{
			"Uncompressed string",
			fields{
				1.0,
				[]string{"caso", "cena", "delfino"},
			},
			args{
				"ca",
			},
			[]string{"caso"},
			false,
		},
		{
			"Compressed string",
			fields{
				1.0,
				[]string{"caso", "cena", "delfino"},
			},
			args{
				"ce",
			},
			[]string{"cena"},
			false,
		},
		{
			"Two strings",
			fields{
				1.0,
				[]string{"caso", "cat", "cena", "delfino"},
			},
			args{
				"ca",
			},
			[]string{"caso", "cat"},
			false,
		},
		{
			"Compressed string next to a compressed one",
			fields{
				1.0,
				[]string{"caso", "cat", "cena", "delfino"},
			},
			args{
				"no",
			},
			[]string{},
			false,
		},
		{
			"Last string",
			fields{
				1.0,
				[]string{"caso", "cat", "cena", "delfino"},
			},
			args{
				"de",
			},
			[]string{"delfino"},
			false,
		},
		{
			"Last string",
			fields{
				1.0,
				[]string{"caso", "cat", "cena", "delfino"},
			},
			args{
				"delfino",
			},
			[]string{"delfino"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lprc := NewLPRC(tt.fields.strings, tt.fields.Epsilon)
			for i, s := range tt.fields.strings {
				lprc.add(s, uint64(i))
			}
			got, err := lprc.FullPrefixSearch(tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("LPRC.FullPrefixSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LPRC.FullPrefixSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}
