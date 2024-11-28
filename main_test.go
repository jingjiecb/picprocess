package main

import (
	"testing"
)

func Test_compressFileToJPEG(t *testing.T) {
	type args struct {
		filePath string
		maxKB    int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case, jpg file no need to compress",
			args: args{
				filePath: "./testdata/avatar.jpg",
				maxKB:    30,
			},
			wantErr: false,
		},
		{
			name: "happy case, jpg file with need to compress",
			args: args{
				filePath: "./testdata/avatar.jpg",
				maxKB:    5,
			},
			wantErr: false,
		},
		{
			name: "happy case, png file with need to compress to jpg",
			args: args{
				filePath: "./testdata/github.png",
				maxKB:    10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compressFileToJPEG(tt.args.filePath, tt.args.maxKB)
			if (err != nil) != tt.wantErr {
				t.Errorf("compressFileToJPEG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.IsSizeLessOrEqualThan(tt.args.maxKB) {
				t.Errorf("compressFileToJPEG() = %v, want less or equal than %v", got, tt.args.maxKB)
			}
		})
	}
}
