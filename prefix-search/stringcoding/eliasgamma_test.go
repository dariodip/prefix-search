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
		name string
		args args
		want uint64
	}{
		{
			"eight bit",
			args{[]string{"a"}},
			uint64(7), // 2 * log_2(8) + 1 = 2 * 3 + 1 = 6 + 1 = 7
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEliasGammaLength(tt.args.strings); got != tt.want {
				t.Errorf("getEliasGammaLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoding_encodeEliasGamma(t *testing.T) {
	type fields struct {
		Strings          *bd.BitData
		Starts           *bd.BitData
		Lengths          *bd.BitData
		LastString       *bd.BitData
		NextIndex        uint64
		NextLengthsIndex uint64
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
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Coding{
				Strings:          tt.fields.Strings,
				Starts:           tt.fields.Starts,
				Lengths:          tt.fields.Lengths,
				LastString:       tt.fields.LastString,
				NextIndex:        tt.fields.NextIndex,
				NextLengthsIndex: tt.fields.NextLengthsIndex,
			}
			if err := c.encodeEliasGamma(tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("Coding.encodeEliasGamma() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCoding_decodeIthEliasGamma(t *testing.T) {
	type fields struct {
		Strings          *bd.BitData
		Starts           *bd.BitData
		Lengths          *bd.BitData
		LastString       *bd.BitData
		NextIndex        uint64
		NextLengthsIndex uint64
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
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Coding{
				Strings:          tt.fields.Strings,
				Starts:           tt.fields.Starts,
				Lengths:          tt.fields.Lengths,
				LastString:       tt.fields.LastString,
				NextIndex:        tt.fields.NextIndex,
				NextLengthsIndex: tt.fields.NextLengthsIndex,
			}
			got, err := c.decodeIthEliasGamma(tt.args.u)
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
		Strings          *bd.BitData
		Starts           *bd.BitData
		Lengths          *bd.BitData
		LastString       *bd.BitData
		NextIndex        uint64
		NextLengthsIndex uint64
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
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Coding{
				Strings:          tt.fields.Strings,
				Starts:           tt.fields.Starts,
				Lengths:          tt.fields.Lengths,
				LastString:       tt.fields.LastString,
				NextIndex:        tt.fields.NextIndex,
				NextLengthsIndex: tt.fields.NextLengthsIndex,
			}
			got, err := c.extractNumFromBinary(tt.args.currentIndex, tt.args.zeroCount)
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
		Strings          *bd.BitData
		Starts           *bd.BitData
		Lengths          *bd.BitData
		LastString       *bd.BitData
		NextIndex        uint64
		NextLengthsIndex uint64
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
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Coding{
				Strings:          tt.fields.Strings,
				Starts:           tt.fields.Starts,
				Lengths:          tt.fields.Lengths,
				LastString:       tt.fields.LastString,
				NextIndex:        tt.fields.NextIndex,
				NextLengthsIndex: tt.fields.NextLengthsIndex,
			}
			got, err := c.eliasGammaZeroCount(tt.args.idx)
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

func TestCoding_eliasGammaZeroCountLoop(t *testing.T) {
	type fields struct {
		Strings          *bd.BitData
		Starts           *bd.BitData
		Lengths          *bd.BitData
		LastString       *bd.BitData
		NextIndex        uint64
		NextLengthsIndex uint64
	}
	type args struct {
		idx       uint64
		zeroCount uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint64
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Coding{
				Strings:          tt.fields.Strings,
				Starts:           tt.fields.Starts,
				Lengths:          tt.fields.Lengths,
				LastString:       tt.fields.LastString,
				NextIndex:        tt.fields.NextIndex,
				NextLengthsIndex: tt.fields.NextLengthsIndex,
			}
			got, err := c.eliasGammaZeroCountLoop(tt.args.idx, tt.args.zeroCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Coding.eliasGammaZeroCountLoop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Coding.eliasGammaZeroCountLoop() = %v, want %v", got, tt.want)
			}
		})
	}
}
