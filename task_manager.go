package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"sync"
	"time"
)

type URLStatus struct {
	URL           string    `json:"url"`
	VisitCount    int       `json:"visitCount"`
	LastVisitTime time.Time `json:"lastVisitTime"`
}

type Task struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"createdAt"`
	MinInterval int         `json:"minInterval"`
	MaxInterval int         `json:"maxInterval"`
	Timeout     int         `json:"timeout"`
	URLs        []string    `json:"urls"`
	URLStatuses []URLStatus `json:"urlStatuses"`
	stopChan    chan bool
	mutex       sync.RWMutex
}

type TaskManager struct {
	tasks map[string]*Task
	mutex sync.RWMutex
}

type CreateTaskRequest struct {
	Name        string   `json:"name" binding:"required"`
	URLs        []string `json:"urls" binding:"required"`
	MinInterval int      `json:"minInterval"`
	MaxInterval int      `json:"maxInterval"`
	Timeout     int      `json:"timeout"`
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*Task),
	}
}

func (tm *TaskManager) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := &Task{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Status:      "running",
		CreatedAt:   time.Now(),
		MinInterval: req.MinInterval,
		MaxInterval: req.MaxInterval,
		Timeout:     req.Timeout,
		URLs:        req.URLs,
		URLStatuses: make([]URLStatus, len(req.URLs)),
		stopChan:    make(chan bool),
	}

	// Initialize URL statuses
	for i, url := range req.URLs {
		task.URLStatuses[i] = URLStatus{
			URL:           url,
			VisitCount:    0,
			LastVisitTime: time.Time{},
		}
	}

	tm.mutex.Lock()
	tm.tasks[task.ID] = task
	tm.mutex.Unlock()

	// Start task execution
	go tm.runTask(task)

	c.JSON(http.StatusOK, task)
}

func (tm *TaskManager) GetTasks(c *gin.Context) {
	tm.mutex.RLock()
	tasks := make([]*Task, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		tasks = append(tasks, task)
	}
	tm.mutex.RUnlock()

	c.JSON(http.StatusOK, tasks)
}

func (tm *TaskManager) GetTaskDetails(c *gin.Context) {
	taskID := c.Param("id")

	tm.mutex.RLock()
	task, exists := tm.tasks[taskID]
	tm.mutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	task.mutex.RLock()
	urlStatuses := task.URLStatuses
	task.mutex.RUnlock()

	c.JSON(http.StatusOK, urlStatuses)
}

func (tm *TaskManager) StopTask(c *gin.Context) {
	taskID := c.Param("id")

	tm.mutex.RLock()
	task, exists := tm.tasks[taskID]
	tm.mutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	task.mutex.Lock()
	if task.Status == "running" {
		task.Status = "stopped"
		close(task.stopChan)
	}
	task.mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Task stopped successfully"})
}

func (tm *TaskManager) runTask(task *Task) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for i := range task.URLs {
		go func(index int) {
			for {
				select {
				case <-task.stopChan:
					return
				default:
					url := task.URLs[index]
					resp, err := client.Get(url)
					
					task.mutex.Lock()
					if err != nil {
						fmt.Printf("Error visiting %s: %v\n", url, err)
					} else {
						io.Copy(io.Discard, resp.Body)
						resp.Body.Close()
						task.URLStatuses[index].VisitCount++
						task.URLStatuses[index].LastVisitTime = time.Now()
					}
					task.mutex.Unlock()

					// Random sleep between MinInterval and MaxInterval
					sleepTime := time.Duration(task.MinInterval+
						time.Now().Nanosecond()%(task.MaxInterval-task.MinInterval+1)) * time.Second
					time.Sleep(sleepTime)
				}
			}
		}(i)
	}

	// Set up timeout
	if task.Timeout > 0 {
		time.Sleep(time.Duration(task.Timeout) * time.Second)
		task.mutex.Lock()
		if task.Status == "running" {
			task.Status = "completed"
			close(task.stopChan)
		}
		task.mutex.Unlock()
	}
}
