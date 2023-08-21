package main

import "testing"

func Test_getTypeId(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		args   args
		wantT  string
		wantId string
	}{
		{
			name:   "test1",
			args:   args{"https://inference-6370f162fe343ed962fe74e7.pool1.mindspore.cn/"},
			wantT:  "inference",
			wantId: "6370f162fe343ed962fe74e7",
		},
		{
			name:   "test2",
			args:   args{"https://evaluate-6370a37ad951790a0933efa7.pool1.mindspore.cn/"},
			wantT:  "evaluate",
			wantId: "6370a37ad951790a0933efa7",
		},
		{
			name:   "test3",
			args:   args{"https://cloud-0c0978e7-bbfb-4182-9582-c3a3bd9e20b6.pool1.mindspore.cn/?token=f9a3bd4e9f2c3be01cd629154cfb224c2703181e050254b5"},
			wantT:  "cloud",
			wantId: "0c0978e7-bbfb-4182-9582-c3a3bd9e20b6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, gotId := getTypeId(tt.args.url)
			if gotT != tt.wantT {
				t.Errorf("getTypeId() gotT = %v, want %v", gotT, tt.wantT)
			}
			if gotId != tt.wantId {
				t.Errorf("getTypeId() gotId = %v, want %v", gotId, tt.wantId)
			}
		})
	}
}
