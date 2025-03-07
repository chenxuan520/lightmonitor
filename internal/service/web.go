package service

import (
	"log"
	"net/http"
	"time"

	"github.com/chenxuan520/lightmonitor/internal/cron"
	"github.com/chenxuan520/lightmonitor/internal/monitor"
	"github.com/chenxuan520/lightmonitor/internal/notify"
	"github.com/gin-gonic/gin"
)

type Web struct {
	MonitorCron *monitor.Cron
	Cron        *cron.Cron
}

func NewWeb() *Web {
	m := monitor.NewCron()
	go m.Run()

	c := cron.NewCron(log.Default())
	go c.Run()

	return &Web{
		MonitorCron: m,
		Cron:        c,
	}
}

// Response 返回值
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	ErrHint string      `json:"err_hint,omitempty"`
}

// Success 成功
func Success(g *gin.Context, data interface{}) {
	g.JSON(http.StatusOK, Response{
		Code: 200,
		Data: data,
	})
}

// Error 错误
func Error(g *gin.Context, status int, data string) {
	g.JSON(status, Response{
		Code: status,
		Data: data,
	})
}

type NotifyMsgReq struct {
	Msg           notify.NotifyMsg `json:"msg"`
	NotifyTime    int64            `json:"notify_time"`
	Notifications []string         `json:"notifications"`
}

func (w *Web) ConfirmTask(g *gin.Context) {
	log.Println("INFO: Confirm received.")
	if len(w.MonitorCron.Tasks) == 0 {
		Success(g, "ok")
		return
	}
	w.MonitorCron.ConfirmChan <- struct{}{}
	Success(g, "ok")
}

func (w *Web) ListTasks(g *gin.Context) {
	log.Println("INFO: List received.")
	if len(w.MonitorCron.Tasks) == 0 {
		Success(g, []monitor.CronTask{})
		return
	}
	result := make(chan struct{ Tasks []monitor.CronTask })
	w.MonitorCron.SnapshotChan <- result
	tasks := <-result
	Success(g, tasks.Tasks)
}

func (w *Web) NotifyMsg(g *gin.Context) {
	log.Println("INFO: NotifyMsg received.")
	var req NotifyMsgReq
	if err := g.BindJSON(&req); err != nil {
		Error(g, http.StatusBadRequest, err.Error())
		return
	}
	if req.NotifyTime == 0 {
		for _, n := range req.Notifications {
			err := notify.SendNotify(n, req.Msg)
			if err != nil {
				Error(g, http.StatusInternalServerError, err.Error())
				return
			}
		}
	} else {
		task := &cron.NotifyOnceTask{
			TaskName: "notify_once_task_" + req.Msg.Title + time.Now().String(),
			NotifyMsg: notify.NotifyMsg{
				Title:   req.Msg.Title,
				Content: req.Msg.Content,
			},
			RunTime:    req.NotifyTime,
			NotifyWays: req.Notifications,
		}
		w.Cron.AddTask(task)

	}
	Success(g, "ok")
}
