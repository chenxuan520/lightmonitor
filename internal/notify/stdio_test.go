package notify

import "testing"

func TestStdio_Send(t *testing.T) {
	type fields struct {
		AbstractNotify AbstractNotify
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
			s := &Stdio{
				AbstractNotify: tt.fields.AbstractNotify,
			}
			if err := s.Send(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Stdio.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
