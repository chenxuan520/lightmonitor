package notify

import (
	"testing"

	"github.com/chenxuan520/lightmonitor/internal/config"
)

func TestFeishu_Send(t *testing.T) {
	config.InitWithPath("../../config/config.json")

	type fields struct {
		AbstractNotify AbstractNotify
		WebHook        string
	}
	type args struct {
		msg NotifyMsg
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"test0",
			fields{
				AbstractNotify: AbstractNotify{},
				WebHook:        "",
			},
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
			f := NewFeishu()
			if err := f.Send(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Feishu.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
