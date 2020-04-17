package types

import (
	"testing"
)

func TestSecureCompare(t *testing.T) {
	type args struct {
		given  []byte
		actual []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "not equal, same size",
			args: args{
				given:  []byte{0x01},
				actual: []byte{0x02},
			},
			want: false,
		},
		{
			name: "not equal, different size",
			args: args{
				given:  []byte{0x01, 0x02},
				actual: []byte{0x02},
			},
			want: false,
		},
		{
			name: "equal, different size",
			args: args{
				given:  []byte{0x00},
				actual: []byte{},
			},
			want: false,
		},
		{
			name: "equal, same size",
			args: args{
				given:  []byte{0x01},
				actual: []byte{0x01},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SecureCompare(tt.args.given, tt.args.actual); got != tt.want {
				t.Errorf("SecureCompare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecureCompareString(t *testing.T) {
	type args struct {
		given  string
		actual string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "not equal, same size",
			args: args{
				given:  "a",
				actual: "b",
			},
			want: false,
		},
		{
			name: "not equal, different size",
			args: args{
				given:  "ab",
				actual: "a",
			},
			want: false,
		},
		{
			name: "equal, different size",
			args: args{
				given:  "\x00",
				actual: "",
			},
			want: false,
		},
		{
			name: "equal, same size",
			args: args{
				given:  "a",
				actual: "a",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SecureCompareString(tt.args.given, tt.args.actual); got != tt.want {
				t.Errorf("SecureCompareString() = %v, want %v", got, tt.want)
			}
		})
	}
}
