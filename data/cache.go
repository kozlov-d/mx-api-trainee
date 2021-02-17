package data

import (
	"fmt"
	"sync"
	"time"

	"github.com/kozlov-d/mx-api-trainee/common"
)

//Cache stores tasks and provides safe access for RW-actions
type Cache struct {
	tasks map[int]*common.Task
	sync.RWMutex
}

//CreateCache return cache object
//Maybe a better way is to init map with size > 0?
func CreateCache() *Cache {
	return &Cache{tasks: make(map[int]*common.Task, 0)}
}

//CreateTask returns task ID
func (c *Cache) CreateTask() int {
	c.RLock()
	ID := len(c.tasks) + 1
	c.RUnlock()

	c.Lock()
	c.tasks[ID] = &common.Task{}
	c.Unlock()

	return ID
}

//GetTaskByID returns task by value
func (c *Cache) GetTaskByID(ID int) (common.Task, error) {
	c.RLock()
	task, found := c.tasks[ID]
	c.RUnlock()
	if found != true {
		return common.Task{}, fmt.Errorf("Task not found with ID=%d", ID)
	}
	return *task, nil
}

//UpdateTask respectively adds given task fields to task fields in cache
func (c *Cache) UpdateTask(key int, t common.Task) {
	c.Lock()
	val := c.tasks[key]
	val.Created += t.Created
	val.Updated += t.Updated
	val.Missed += t.Missed
	val.Deleted += t.Deleted
	c.tasks[key] = val
	c.Unlock()
}

//CompleteTask sets the flag of completion to 1
func (c *Cache) CompleteTask(ID int, start time.Time) {
	c.Lock()
	c.tasks[ID].IsCompleted = true
	c.tasks[ID].TimeSpent = fmt.Sprintf("%s", time.Since(start))
	c.Unlock()
}
