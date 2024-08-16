package edgeTTS

import (
	"testing"
)

func Test_listVoices(t *testing.T) {
	tests := []struct {
		name    string
		want    []Voice
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "test-1",
			want:    []Voice{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listVoices()
			if len(got) <= 0 {
				t.Errorf("listVoices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestVoicesManager_find(t *testing.T) {
	type args struct {
		attributes Voice
	}
	vm := &VoicesManager{}
	_ = vm.create(nil)
	tests := []struct {
		name string
		vm   *VoicesManager
		args args
		want []Voice
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			vm:   vm,
			args: args{
				attributes: Voice{
					Locale: "zh-CN",
				},
			},
			want: []Voice{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.vm.find(tt.args.attributes)
			if len(got) <= 0 {
				t.Errorf("listVoices() wantErr %v", tt.want)
				return
			}
			t.Logf("listVoices() got %v", got)
		})
	}
}
