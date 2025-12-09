package utils

import (
	"fmt"
	"time"
)

// Scheduler manages scheduled tasks
type Scheduler struct {
	tasks []Task
}

// Task represents a scheduled task
type Task struct {
	Name     string
	Interval time.Duration
	Function func()
	LastRun  time.Time
}

// NewScheduler creates a new scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		tasks: make([]Task, 0),
	}
}

// AddTask adds a new task to the scheduler
func (s *Scheduler) AddTask(name string, interval time.Duration, fn func()) {
	s.tasks = append(s.tasks, Task{
		Name:     name,
		Interval: interval,
		Function: fn,
		LastRun:  time.Now(),
	})
}

// Run starts the scheduler
func (s *Scheduler) Run() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		for i, task := range s.tasks {
			if time.Since(task.LastRun) >= task.Interval {
				fmt.Printf("Running scheduled task: %s\n", task.Name)
				go task.Function()
				s.tasks[i].LastRun = time.Now()
			}
		}
	}
}

// CheckMonthlyZakat checks if zakat should be deducted (runs on 1st of month)
func CheckMonthlyZakat() bool {
	now := time.Now()
	return now.Day() == 1 && now.Hour() == 0 && now.Minute() == 0
}
