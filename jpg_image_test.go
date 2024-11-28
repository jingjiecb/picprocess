package main

import (
	"bytes"
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestNewJPGImage(t *testing.T) {
	type args struct {
		image   image.Image
		quality int
	}
	mockImage := image.NewNRGBA(image.Rect(0, 0, 256, 256))

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				image:   mockImage,
				quality: 80,
			},
			wantErr: false,
		},
		{
			name: "invalid quality",
			args: args{
				image:   mockImage,
				quality: -1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewJPGImage(tt.args.image, tt.args.quality)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJPGImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestJPGImage_IsSizeLessThanKB(t *testing.T) {
	type fields struct {
		content *bytes.Buffer
	}
	type args struct {
		sizeKB int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "small image",
			fields: fields{
				content: generateBuffer(100 * 1024),
			},
			args: args{
				sizeKB: 100,
			},
			want: true,
		},
		{
			name: "large image",
			fields: fields{
				content: generateBuffer(301 * 1024),
			},
			args: args{
				sizeKB: 300,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jpg := &JPGImage{
				content: tt.fields.content,
			}
			if got := jpg.IsSizeLessOrEqualThan(tt.args.sizeKB); got != tt.want {
				t.Errorf("IsSizeLessThanKB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func generateBuffer(sizeB int) *bytes.Buffer {
	data := make([]byte, sizeB)

	for i := range data {
		data[i] = 'A'
	}

	return bytes.NewBuffer(data)
}

func TestJPGImage_SaveToFile(t *testing.T) {
	type fields struct {
		content *bytes.Buffer
	}
	mockContent := generateBuffer(1)

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			fields: fields{
				content: mockContent,
			},
			args: args{
				path: filepath.Join(os.TempDir(), "picprocess-test-"+uuid.New().String()+".jpg"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Cleanup(func() {
			err := os.Remove(tt.args.path)
			if err != nil {
				t.Error(err)
			}
		})
		t.Run(tt.name, func(t *testing.T) {
			jpg := &JPGImage{
				content: tt.fields.content,
			}
			if err := jpg.SaveToFile(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			_, err := os.Stat(tt.args.path)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestNewJPGImageFromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				path: "./testdata/avatar.jpg",
			},
			wantErr: false,
		},
		{
			name: "not a jpeg format image",
			args: args{
				path: "./testdata/github.png",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.args.path)
			if err != nil {
				t.Error(err)
			}
			defer file.Close()

			_, err = NewJPGImageFromFile(file)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJPGImageFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
