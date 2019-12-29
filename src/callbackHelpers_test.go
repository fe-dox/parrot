package main

import (
	"testing"
)

func TestDecodeString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name     string
		args     args
		wantCsid CallbackStackItemID
		wantCiid CallbackItemID
		wantErr  bool
	}{
		{
			"Check if 1-1 returns 1 1",
			args{s: "1-1"},
			1,
			1,
			false,
		},
		{
			"Check if error is thrown",
			args{s: "2-asdsiji331"},
			0,
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCsid, gotCiid, err := DecodeString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCsid != tt.wantCsid {
				t.Errorf("DecodeString() gotCsid = %v, want %v", gotCsid, tt.wantCsid)
			}
			if gotCiid != tt.wantCiid {
				t.Errorf("DecodeString() gotCiid = %v, want %v", gotCiid, tt.wantCiid)
			}
		})
	}
}

func TestPrepareString(t *testing.T) {
	type args struct {
		csid CallbackStackItemID
		ciid CallbackItemID
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Check if string is built correctly",
			args{csid: 123, ciid: 123},
			"123-123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prepareString(tt.args.csid, tt.args.ciid); got != tt.want {
				t.Errorf("prepareString() = %v, want %v", got, tt.want)
			}
		})
	}
}
