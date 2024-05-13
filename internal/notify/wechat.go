package notify

import (
	"context"
	"fmt"

	"github.com/chenxuan520/lightmonitor/internal/config"
	utils "github.com/chenxuan520/lightmonitor/internal/utils"
)

type WeChat struct {
	AbstractNotify
	PostUrl string
}

func NewWeChat() *WeChat {
	return &WeChat{
		PostUrl: fmt.Sprintf("https://sctapi.ftqq.com/%s.send", config.GlobalConfig.NotifyWay.Wechat.SendKey),
	}
}

func (w *WeChat) Send(msg NotifyMsg) error {
	if w.PostUrl == "" {
		return fmt.Errorf("wechat not init")
	}
	_, err := utils.HttpRequest(context.Background(), "GET", w.PostUrl, nil, nil, map[string]string{
		"title": msg.Title,
		"desp":  msg.Content,
	})
	return err
}
