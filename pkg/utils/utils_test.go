package utils

import (
	"reflect"
	"testing"
)

func Test_PadBytes(t *testing.T) {
	type args struct {
		inputBytes []byte
		blockSize  int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "successful padding of string to block size",
			args: args{
				inputBytes: []byte("hello"),
				blockSize:  10,
			},
			want:    []byte("hello-----"),
			wantErr: false,
		},
		{
			name: "successful padding of string to block size",
			args: args{
				inputBytes: []byte("hello"),
				blockSize:  5,
			},
			want:    []byte("hello"),
			wantErr: false,
		},
		{
			name: "failed padding of string to block size - input too long",
			args: args{
				inputBytes: []byte("hellohello"),
				blockSize:  8,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PadBytes(tt.args.inputBytes, tt.args.blockSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("padBytesToBlockSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("padBytesToBlockSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBlockSize(t *testing.T) {
	type TestData struct {
		Name    string `maxBytes:"100"`
		Comment string `maxBytes:"1000"`
	}
	type TestDataBad struct {
		Name    string `maxBytes:"100"`
		Comment string
	}
	type TestDataBadTag struct {
		Name    string `maxBytes:"100"`
		Comment string `maxBytes:"1000wrong"`
	}

	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{
			name:    "simple check",
			want:    1100,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBlockSize[TestData]()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBlockSize() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("require maxBytes tags", func(t *testing.T) {
		_, err := GetBlockSize[TestDataBad]()
		if err == nil {
			t.Error("GetBlockSize() expected error got nil")
		}
	})

	t.Run("require maxBytes tags to be int", func(t *testing.T) {
		_, err := GetBlockSize[TestDataBadTag]()
		if err == nil {
			t.Error("GetBlockSize() expected error got nil")
		}
	})
}
