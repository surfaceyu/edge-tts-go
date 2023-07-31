package edgeTTS

import (
	"testing"
)

func Test_getHeadersAndData(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		want1   []byte
		wantErr bool
	}{
		{
			name: "test-1",
			args: args{
				data: "X-Timestamp:2022-01-01\r\nContent-Type:application/json; charset=utf-8\r\nPath:speech.config\r\n\r\n{\"context\":{\"synthesis\":{\"audio\":{\"metadataoptions\":{\"sentenceBoundaryEnabled\":false,\"wordBoundaryEnabled\":true},\"outputFormat\":\"audio-24khz-48kbitrate-mono-mp3\"}}}}",
			},
			want:    map[string]string{},
			want1:   []byte{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getHeadersAndData(tt.args.data)
			t.Logf("%v \n%v \n", got, got1)
			if (err != nil) != tt.wantErr {
				t.Errorf("getHeadersAndData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("getHeadersAndData() got = %v, want %v", got, tt.want)
			// }
			// if !reflect.DeepEqual(got1, tt.want1) {
			// 	t.Errorf("getHeadersAndData() got1 = %v, want %v", got1, tt.want1)
			// }
		})
	}
}
