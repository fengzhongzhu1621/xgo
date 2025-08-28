package handler

import (
	"context"
	"net/http"

	"github.com/fengzhongzhu1621/xgo/cron/scheduler"
	"github.com/gin-gonic/gin"
)

var JobFuncMap = map[string]func(ctx context.Context){
	"demo": func(ctx context.Context) {
		println("Task executed")
	},
}

type TaskRequestSerializer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
}

func AddTask(c *gin.Context) {
	var req TaskRequestSerializer

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获得任务执行函数
	jobFunc := JobFuncMap[req.ID]
	if jobFunc == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "job func not found"})
		return
	}

	// 添加任务
	task := &scheduler.Task{
		ID:       req.ID,
		Name:     req.Name,
		Schedule: req.Schedule,
		JobFunc:  jobFunc,
	}
	sched := scheduler.GetScheduler()
	err := sched.AddTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task added"})
}

func RunNow(c *gin.Context) {
	id := c.Param("id")

	sched := scheduler.GetScheduler()
	sched.RunNow(id)

	c.JSON(http.StatusOK, gin.H{"message": "task executed now"})
}

func PauseTask(c *gin.Context) {
	id := c.Param("id")

	sched := scheduler.GetScheduler()

	sched.PauseTask(id)

	c.JSON(http.StatusOK, gin.H{"message": "task paused"})
}

func ResumeTask(c *gin.Context) {
	id := c.Param("id")

	sched := scheduler.GetScheduler()

	err := sched.ResumeTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task resumed"})
}
