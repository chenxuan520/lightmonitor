package cron

import (
	"log"
	"testing"
	"time"

	"github.com/chenxuan520/lightmonitor/internal/notify"
)

func TestCron_Run(t *testing.T) {
	t.Log("TestCron_Run")
	c := NewCron(log.Default())
	c.AddTask(&NotifyOnceTask{
		TaskName: "test",
		NotifyMsg: notify.NotifyMsg{
			Title:   "test_title",
			Content: "test_content",
		},
		RunTime:    time.Now().Unix() + 16,
		NotifyWays: []string{"stdio"},
	})
	c.AddTask(&NotifyOnceTask{
		TaskName: "1",
		NotifyMsg: notify.NotifyMsg{
			Title:   "test_title",
			Content: "test_content",
		},
		RunTime:    time.Now().Unix() + 3,
		NotifyWays: []string{"stdio"},
	})
	go func() {
		err := c.Run()
		if err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(4 * time.Second)
	c.Stop()
}
