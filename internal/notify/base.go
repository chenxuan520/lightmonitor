package notify

import "fmt"

type NotifyMsg struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Notify interface {
	Send(msg NotifyMsg) error
}

type AbstractNotify struct {
}

func (a *AbstractNotify) Send(msg NotifyMsg) error {
	return fmt.Errorf("Not implemented")
}

// 策略模式
func SendNotify(notify string, msg NotifyMsg) error {
	n := GetNotifyByStr(notify)
	if n == nil {
		return fmt.Errorf("notify not found")
	}
	return n.Send(msg)
}

func GetNotifyByStr(notify string) Notify {
	switch notify {
	case "email":
		return NewEmail()
	case "wechat":
		return NewWeChat()
	case "feishu":
		return NewFeishu()
	case "stdio":
		return NewStdio()
	default:
		return nil
	}
}
