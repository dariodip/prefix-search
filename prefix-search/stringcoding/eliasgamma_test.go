package stringcoding

import (
	"testing"

	bd "github.com/dariodip/prefix-search/prefix-search/bitdata"
)

func Test_getEliasGammaLength(t *testing.T) {
	type args struct {
		strings []string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{
			"eight bit",
			args{[]string{"a"}},
			uint64(7), // 2 * log_2(8) + 1 = 2 * 3 + 1 = 6 + 1 = 7
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getEliasGammaLength(tt.args.strings);
			if (err != nil) != tt.wantErr {
				t.Errorf("Coding.getEliasGammaLength() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getEliasGammaLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoding_encodeEliasGamma(t *testing.T) {
	type fields struct {
		// stub values
		strings []string
		epsilon float64
	}
	type args struct {
		n uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"string 1 length coding",
			fields{[]string{"stub"}, 20},
			args{24},
			false,
		},
		{
			"empty string",
			fields{[]string{"stub"}, 20},
			args{0},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lprc := NewLPRC(tt.fields.strings, tt.fields.epsilon)
			if err := lprc.coding.encodeEliasGamma(tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("Coding.encodeEliasGamma() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCoding_decodeIthEliasGamma(t *testing.T) {
	type fields struct {
		// stub values
		strings []string
		epsilon float64
	}
	type args struct {
		u uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint64
		wantErr bool
	}{
		{
			"first string length",
			fields{[]string{"ciao", "cic"}, 20},
			args{0},
			bd.GetLengthInBit("ciao"),
			false,
		},
		{
			"second string length",
			fields{[]string{"ciao", "cic"}, 20},
			args{1},
			uint64(2),
			false,
		},
		{
			"error test",
			fields{[]string{"ciao", "cic"}, 20},
			args{3},
			uint64(0),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lprc := NewLPRC(tt.fields.strings, tt.fields.epsilon)
			for index, s := range tt.fields.strings {
				lprc.add(s, uint64(index))
			}
			got, err := lprc.coding.decodeIthEliasGamma(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("Coding.decodeIthEliasGamma() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Coding.decodeIthEliasGamma() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoding_extractNumFromBinary(t *testing.T) {
	type fields struct {
		// stub values
		strings []string
		epsilon float64
	}
	type args struct {
		currentIndex uint64
		zeroCount    uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint64
		wantErr bool
	}{
		{
			"first string length",
			fields{[]string{"ciao", "cic"}, 20},
			args{5, 5},
			bd.GetLengthInBit("ciao"),
			false,
		},
		{
			"second string length",
			fields{[]string{"ciao", "cic"}, 20},
			args{12, 1},
			uint64(2),
			false,
		},
		{
			"error test",
			fields{[]string{"ciao", "cic"}, 20},
			args{15, 0},
			uint64(0),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lprc := NewLPRC(tt.fields.strings, tt.fields.epsilon)
			for index, s := range tt.fields.strings {
				lprc.add(s, uint64(index))
			}
			got, err := lprc.coding.extractNumFromBinary(tt.args.currentIndex, tt.args.zeroCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Coding.extractNumFromBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Coding.extractNumFromBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoding_eliasGammaZeroCount(t *testing.T) {
	type fields struct {
		// stub values
		strings []string
		epsilon float64
	}
	type args struct {
		idx uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint64
		wantErr bool
	}{
		{
			"first string length",
			fields{[]string{"ciao", "cic"}, 20},
			args{0},
			uint64(5),
			false,
		},
		{
			"second string length",
			fields{[]string{"ciao", "cic"}, 20},
			args{11},
			uint64(1),
			false,
		},
		{
			"error test",
			fields{[]string{"ciao", "cic"}, 20},
			args{15},
			uint64(0),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lprc := NewLPRC(tt.fields.strings, tt.fields.epsilon)
			for index, s := range tt.fields.strings {
				lprc.add(s, uint64(index))
			}
			got, err := lprc.coding.eliasGammaZeroCount(tt.args.idx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Coding.eliasGammaZeroCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Coding.eliasGammaZeroCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
