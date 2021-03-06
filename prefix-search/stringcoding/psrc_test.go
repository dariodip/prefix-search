package stringcoding

import (
	"reflect"
	"testing"
)

func TestPSRC_Retrieval(t *testing.T) {
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
			"Last string",
			fields{
				1.0,
				[]string{"caso", "cena", "delfino"},
			},
			args{
				uint64(2),
				uint64(16),
			},
			"de",
			false,
		},
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
			"",
			true,
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
			psrc := NewPSRC(tt.fields.strings, tt.fields.Epsilon)
			for i, s := range tt.fields.strings {
				psrc.add(s, uint64(i))
			}
			got, err := psrc.Retrieval(tt.args.u, tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("PSRC.Retrieval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("len(got)= %d", len(got))
				t.Errorf("PSRC.Retrieval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPSRC_getStringLength(t *testing.T) {
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
			uint64(48),
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
			uint64(48),
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
			uint64(72),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			psrc := NewPSRC(tt.fields.strings, tt.fields.Epsilon)
			for i, s := range tt.fields.strings {
				psrc.add(s, uint64(i))
			}
			got, err := psrc.getStringLength(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: PSRC.getStringLength() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("%s: PSRC.getStringLength() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestPSRC_FullPrefixSearch(t *testing.T) {
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
			"Full string",
			fields{
				1,
				[]string{"casotto", "cisonostatierrori", "cuz", "delfino"},
			},
			args{
				"delfino",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			psrc := NewPSRC(tt.fields.strings, tt.fields.Epsilon)
			for i, s := range tt.fields.strings {
				psrc.add(s, uint64(i))
			}
			got, err := psrc.FullPrefixSearch(tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("Test: %s", tt.name)
				t.Errorf("PSRC.FullPrefixSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s: PSRC.FullPrefixSearch() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
