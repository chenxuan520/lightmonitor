package notify

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chenxuan520/lightmonitor/internal/config"
	apis "github.com/chenxuan520/lightmonitor/internal/utils"
)

type Feishu struct {
	AbstractNotify
	WebHook string
}

type feishuMsg struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

func NewFeishu() *Feishu {
	return &Feishu{
		WebHook: config.GlobalConfig.NotifyWay.Feishu.WebHook,
	}
}

func (f *Feishu) Send(msg NotifyMsg) error {
	if f.WebHook == "" {
		return fmt.Errorf("feishu webhook is empty")
	}
	// build msg
	feishuMsg := feishuMsg{
		MsgType: "text",
		Content: struct {
			Text string `json:"text"`
		}{
			Text: msg.Title + ":" + msg.Content,
		},
	}
	body, err := json.Marshal(feishuMsg)
	if err != nil {
		return err
	}
	// send msg
	apis.HttpRequest(context.Background(), "POST", f.WebHook, []byte(body), map[string]string{
		"Content-Type": "application/json",
	}, nil)
	return nil
}
