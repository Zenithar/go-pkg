package types

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type innerStruct struct {
	String string `commented:"true"`
	Long   uint64
}

type invalidInnerStruct struct {
	String string `commented:"123"`
	Long   uint64
}

func TestAsEnvVariables(t *testing.T) {
	type args struct {
		o             interface{}
		prefix        string
		skipCommented bool
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "nil",
			args: args{
				o:             nil,
				prefix:        "",
				skipCommented: false,
			},
			wantErr: true,
		},
		{
			name: "empty",
			args: args{
				o:             struct{}{},
				prefix:        "",
				skipCommented: false,
			},
			wantErr: false,
			want:    map[string]string{},
		},
		{
			name: "simple",
			args: args{
				o: struct {
					Bool     bool
					String   string
					Duration time.Duration
				}{
					Bool:     true,
					String:   "test",
					Duration: 2 * time.Hour,
				},
				prefix:        "SIMPLE",
				skipCommented: false,
			},
			wantErr: false,
			want: map[string]string{
				"SIMPLE_BOOL":     "true",
				"SIMPLE_DURATION": "2h0m0s",
				"SIMPLE_STRING":   "test",
			},
		},
		{
			name: "simple w/ commented",
			args: args{
				o: struct {
					Bool     bool
					String   string `commented:"true"`
					Duration time.Duration
				}{
					Bool:     true,
					String:   "test",
					Duration: 2 * time.Hour,
				},
				prefix:        "SIMPLE",
				skipCommented: true,
			},
			wantErr: false,
			want: map[string]string{
				"SIMPLE_BOOL":     "true",
				"SIMPLE_DURATION": "2h0m0s",
			},
		},
		{
			name: "simple invalid bool",
			args: args{
				o: struct {
					Bool     bool
					String   string `commented:"123"`
					Duration time.Duration
				}{
					Bool:     true,
					String:   "test",
					Duration: 2 * time.Hour,
				},
				prefix:        "",
				skipCommented: true,
			},
			wantErr: true,
		},
		{
			name: "complex",
			args: args{
				o: struct {
					Bool     bool
					String   string
					Duration time.Duration
					Inner    *innerStruct
				}{
					Bool:     true,
					String:   "test",
					Duration: 2 * time.Hour,
					Inner: &innerStruct{
						Long: 64,
					},
				},
				prefix:        "SIMPLE",
				skipCommented: false,
			},
			wantErr: false,
			want: map[string]string{
				"SIMPLE_BOOL":         "true",
				"SIMPLE_DURATION":     "2h0m0s",
				"SIMPLE_STRING":       "test",
				"SIMPLE_INNER_STRING": "",
				"SIMPLE_INNER_LONG":   "64",
			},
		},
		{
			name: "complex w/ commented",
			args: args{
				o: struct {
					Bool     bool
					String   string
					Duration time.Duration
					Inner    *innerStruct
				}{
					Bool:     true,
					String:   "test",
					Duration: 2 * time.Hour,
					Inner: &innerStruct{
						Long: 64,
					},
				},
				prefix:        "SIMPLE",
				skipCommented: true,
			},
			wantErr: false,
			want: map[string]string{
				"SIMPLE_BOOL":       "true",
				"SIMPLE_DURATION":   "2h0m0s",
				"SIMPLE_STRING":     "test",
				"SIMPLE_INNER_LONG": "64",
			},
		},
		{
			name: "complex w/ commented error",
			args: args{
				o: struct {
					Bool     bool
					String   string
					Duration time.Duration
					Inner    *invalidInnerStruct
				}{
					Bool:     true,
					String:   "test",
					Duration: 2 * time.Hour,
					Inner: &invalidInnerStruct{
						Long: 64,
					},
				},
				prefix:        "SIMPLE",
				skipCommented: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AsEnvVariables(tt.args.o, tt.args.prefix, tt.args.skipCommented)
			if (err != nil) != tt.wantErr {
				t.Errorf("AsEnvVariables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("AsEnvVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}
