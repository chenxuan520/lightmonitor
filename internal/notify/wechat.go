package notify

import (
	"context"
	"fmt"

	"github.com/chenxuan520/lightmonitor/internal/config"
	utils "github.com/chenxuan520/lightmonitor/internal/utils"
)

type WeChat struct {
	AbstractNotify
	PostUrls []string
}

func NewWeChat() *WeChat {
	postUrls := make([]string, len(config.GlobalConfig.NotifyWay.Wechat.SendKeys))
	for _, v := range config.GlobalConfig.NotifyWay.Wechat.SendKeys {
		postUrls = append(postUrls, fmt.Sprintf("https://sctapi.ftqq.com/%s.send", v))

	}
	return &WeChat{
		PostUrls: postUrls,
	}
}

func (w *WeChat) Send(msg NotifyMsg) error {
	if len(w.PostUrls) == 0 {
		return fmt.Errorf("wechat not init")
	}

	for _, v := range w.PostUrls {
		_, err := utils.HttpRequest(context.Background(), "GET", v, nil, nil, map[string]string{
			"title": msg.Title,
			"desp":  msg.Content,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
