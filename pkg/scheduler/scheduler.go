package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/libesz/poolmanager/pkg/controller"
)

func New() Scheduler {
	return Scheduler{
		taskChan:    make(chan schedulerTask),
		cancelChan:  make(chan cancelTask),
		controllers: make(map[string]controller.Controller),
		queue:       make(map[string]chan struct{}),
	}
}

func (s *Scheduler) AddController(c controller.Controller) {
	log.Printf("Scheduler: added controller: %s\n", c.GetName())
	s.controllers[c.GetName()] = c
}

func (s *Scheduler) GetConfigProperties(controllerName string) controller.ConfigProperties {
	c, ok := s.controllers[controllerName]
	if !ok {
		return controller.EmptyProperties()
	}
	return c.GetConfigProperties()
}

func (s *Scheduler) ConfigUpdated(controllerName string, config controller.Config) error {
	c, ok := s.controllers[controllerName]
	if !ok {
		return fmt.Errorf("Controller not found: %s", controllerName)
	}
	if err := c.ValidateConfig(config); err != nil {
		return err
	}
	log.Printf("Scheduler: scheduling controller: %s\n", controllerName)
	s.cancel(controllerName)
	s.enqueue(schedulerTask{controller: c, config: config})
	return nil
}

func (s *Scheduler) cancel(controllerName string) {
	task := cancelTask{controller: controllerName, result: make(chan struct{})}
	s.cancelChan <- task
	<-task.result
}

func (s *Scheduler) enqueue(task schedulerTask) {
	s.taskChan <- task
}

func (s *Scheduler) Run(stopChan chan struct{}) {
	for {
		select {
		case task := <-s.taskChan:
			log.Printf("Scheduler: executing controller: %s\n", task.controller.GetName())
			reEnqueAfterSet := task.controller.Act(task.config)
			for _, reEnqueAfter := range reEnqueAfterSet {
				queueItem := make(chan struct{})
				s.queue[reEnqueAfter.Controller.GetName()] = queueItem
				go func(request controller.EnqueueRequest) {
					timer := time.After(request.After)
					select {
					case <-timer:
						s.enqueue(schedulerTask{controller: request.Controller, config: request.Config})
					case <-queueItem:
						log.Printf("Cancelling task for controller: %s\n", request.Controller.GetName())
					}
				}(reEnqueAfter)
			}
		case cancelRequest := <-s.cancelChan:
			queueItem, ok := s.queue[cancelRequest.controller]
			if ok {
				close(queueItem)
			} else {
				log.Printf("Scheduler: could not cancel any task for controller: %s\n", cancelRequest.controller)
			}
			s.queue[cancelRequest.controller] = nil
			close(cancelRequest.result)
		case <-stopChan:
			return
		}
	}
}
