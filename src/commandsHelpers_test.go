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
			"Check what happens",
			[]string{},
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
