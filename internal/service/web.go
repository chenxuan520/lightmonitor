package service

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Web struct {
	Cron *Cron
}

func NewWeb(c *Cron) *Web {
	return &Web{
		Cron: c,
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

func (w *Web) ConfirmTask(g *gin.Context) {
	log.Println("INFO: Confirm received.")
	w.Cron.ConfirmChan <- struct{}{}
	Success(g, "ok")
}

func (w *Web) ListTasks(g *gin.Context) {
	log.Println("INFO: List received.")
	result := make(chan struct{ Tasks []CronTask })
	w.Cron.SnapshotChan <- result
	tasks := <-result
	Success(g, tasks.Tasks)
}
