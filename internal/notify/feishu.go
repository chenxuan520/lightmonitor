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
	WebHooks []string
}

type feishuMsg struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

func NewFeishu() *Feishu {
	return &Feishu{
		WebHooks: config.GlobalConfig.NotifyWay.Feishu.WebHooks,
	}
}

func (f *Feishu) Send(msg NotifyMsg) error {
	if len(f.WebHooks) == 0 {
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
	for _, webHook := range f.WebHooks {
		_, err := apis.HttpRequest(context.Background(), "POST", webHook, []byte(body), map[string]string{
			"Content-Type": "application/json",
		}, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
