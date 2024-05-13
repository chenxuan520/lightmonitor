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
func SendNotify(notify Notify, msg NotifyMsg) error {
	return notify.Send(msg)
}
