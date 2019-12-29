package main

import (
	"reflect"
	"testing"
)

func Test_listAllDrives(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			"",
			[]string{"C", "D", "E", "F", "X", "Y"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listAllDrives()
			if (err != nil) != tt.wantErr {
				t.Errorf("listAllDrives() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listAllDrives() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name     string
		args     args
		wantDir  string
		wantFile string
	}{
		{
			"",
			args{path: "C:\\Users\\Test\\Test2\\Item.txt"},
			"C:\\Users\\Test\\Test2\\",
			"Item.txt",
		},
		{
			"",
			args{path: "\\\\Server\\Directory\\Item.txt"},
			"\\\\Server\\Directory\\",
			"Item.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDir, gotFile := splitPath(tt.args.path)
			if gotDir != tt.wantDir {
				t.Errorf("splitPath() gotDir = %v, want %v", gotDir, tt.wantDir)
			}
			if gotFile != tt.wantFile {
				t.Errorf("splitPath() gotFile = %v, want %v", gotFile, tt.wantFile)
			}
		})
	}
}
