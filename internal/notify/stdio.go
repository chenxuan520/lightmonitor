package notify

import "fmt"

type Stdio struct {
	AbstractNotify
}

func NewStdio() *Stdio {
	return &Stdio{}
}

func (s *Stdio) Send(msg NotifyMsg) error {
	fmt.Println(msg.Title + ":" + msg.Content)
	return nil
}
