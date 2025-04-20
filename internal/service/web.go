package service

import (
	"log"
	"net/http"
	"strconv"
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
	CycleTime     int64            `json:"cycle_time"`
	Notifications []string         `json:"notifications"`
}

func (w *Web) ConfirmMonitorTask(g *gin.Context) {
	log.Println("INFO: Confirm received.")
	if len(w.MonitorCron.Tasks) == 0 {
		Success(g, "ok")
		return
	}
	w.MonitorCron.ConfirmChan <- struct{}{}
	Success(g, "ok")
}

func (w *Web) ListMonitorTask(g *gin.Context) {
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

func (w *Web) ListTask(g *gin.Context) {
	log.Println("INFO: List received.")
	if len(w.Cron.Tasks) == 0 {
		Success(g, []cron.CronTask{})
		return
	}
	result := make(chan struct{ Tasks []cron.CronTask })
	w.Cron.SnapshotChan <- result
	tasks := <-result
	type respDate struct {
		TaskName    string `json:"task_name"`
		NextRunTime int64  `json:"next_run_time"`
	}
	data := make([]respDate, len(tasks.Tasks))
	for i, task := range tasks.Tasks {
		data[i] = respDate{
			TaskName:    task.Name(),
			NextRunTime: task.NextRunTime(),
		}
	}
	Success(g, map[string]interface{}{
		"tasks": data,
	})
}

func (w *Web) CancelTask(g *gin.Context) {
	log.Println("INFO: Cancel received.")
	var req struct {
		TaskName string `json:"task_name"`
	}
	if err := g.BindJSON(&req); err != nil {
		Error(g, http.StatusBadRequest, err.Error())
		return
	}
	w.Cron.DeleteTaskChan <- req.TaskName
	Success(g, "ok")
}

func (w *Web) NotifyMsg(g *gin.Context) {
	log.Println("INFO: NotifyMsg received.")
	var req NotifyMsgReq
	if err := g.BindJSON(&req); err != nil {
		Error(g, http.StatusBadRequest, err.Error())
		return
	}

	taskName := ""
	titleSuffix := req.Msg.Title + "_" + strconv.FormatInt(time.Now().Unix(), 10)
	if req.NotifyTime == 0 {
		for _, n := range req.Notifications {
			err := notify.SendNotify(n, req.Msg)
			if err != nil {
				Error(g, http.StatusInternalServerError, err.Error())
				return
			}
		}
	} else if req.CycleTime == 0 {
		task := &cron.NotifyOnceTask{
			TaskName: "notify_once_task_" + titleSuffix,
			NotifyMsg: notify.NotifyMsg{
				Title:   req.Msg.Title,
				Content: req.Msg.Content,
			},
			RunTime:    req.NotifyTime,
			NotifyWays: req.Notifications,
		}
		taskName = task.Name()
		err := w.Cron.AddTask(task)
		if err != nil {
			Error(g, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		task := &cron.NotifyCycleTask{
			TaskName:   "notify_cycle_task_" + titleSuffix,
			NotifyMsg:  notify.NotifyMsg{Title: req.Msg.Title, Content: req.Msg.Content},
			RunTime:    req.NotifyTime,
			CycleTime:  req.CycleTime,
			NotifyWays: req.Notifications,
		}
		taskName = task.Name()
		err := w.Cron.AddTask(task)
		if err != nil {
			Error(g, http.StatusInternalServerError, err.Error())
			return
		}
	}
	Success(g, map[string]interface{}{
		"task_name": taskName,
	})
}
