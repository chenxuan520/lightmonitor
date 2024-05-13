package notify

import (
	"testing"

	"github.com/chenxuan520/lightmonitor/internal/config"
)

func TestWeChat_Send(t *testing.T) {
	config.InitWithPath("../../config/config.json")

	type args struct {
		msg NotifyMsg
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test0",
			args{
				msg: NotifyMsg{
					Title:   "测试标题",
					Content: "测试内容",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWeChat()
			if err := w.Send(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("WeChat.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
