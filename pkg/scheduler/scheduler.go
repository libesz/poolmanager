package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/libesz/poolmanager/pkg/configstore"
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

func (s *Scheduler) SetConfigStore(configStore *configstore.ConfigStore) {
	s.configStore = configStore
}

func (s *Scheduler) AddController(c controller.Controller) {
	log.Printf("Scheduler: added controller: %s\n", c.GetName())
	s.controllers[c.GetName()] = c
	var config controller.Config
	if s.configStore != nil {
		config = s.configStore.Get(c.GetName())
	} else {
		config = controller.EmptyConfig()
	}
	s.enqueue(schedulerTask{controller: c, config: config})
}

func (s *Scheduler) GetConfigProperties(controllerName string) controller.ConfigProperties {
	c, ok := s.controllers[controllerName]
	if !ok {
		return controller.EmptyProperties()
	}
	return c.GetConfigProperties()
}

func (s *Scheduler) ConfigUpdated(controllerName string, config controller.Config, enqueue bool) error {
	c, ok := s.controllers[controllerName]
	if !ok {
		return fmt.Errorf("Scheduler: Controller not found: %s", controllerName)
	}
	if err := c.ValidateConfig(config); err != nil {
		return err
	}
	log.Printf("Scheduler: config changed for controller: %s\n", controllerName)
	if enqueue {
		s.cancelOuter(controllerName)
		s.enqueue(schedulerTask{controller: c, config: config})
	}
	return nil
}

func (s *Scheduler) cancelOuter(controllerName string) {
	task := cancelTask{controller: controllerName, result: make(chan struct{})}
	s.cancelChan <- task
	<-task.result
}

func (s *Scheduler) cancelInner(cancelItemChan chan struct{}) {
	select {
	case cancelItemChan <- struct{}{}:
		<-cancelItemChan
	default:
		// Task is already done, no go-routing is listening on the channel
	}
}

func (s *Scheduler) enqueue(task schedulerTask) {
	s.taskChan <- task
}

func (s *Scheduler) taskInQueue(request controller.EnqueueRequest, cancelChan chan struct{}) {
	if s.configStore != nil && !controller.IsEmptyConfig(request.Config) {
		if err := s.configStore.Set(request.Controller.GetName(), request.Config, false); err != nil {
			log.Printf("Scheduler: invalid config pushed back by controller %s: %s. Aborting controller.\n", request.Controller.GetName(), err.Error())
			return
		}
	}
	timer := time.After(request.After)
	select {
	case <-timer:
		log.Printf("Scheduler: enqueing task for controller: %s\n", request.Controller.GetName())
		s.enqueue(schedulerTask{controller: request.Controller, config: request.Config})
	case <-cancelChan:
		log.Printf("Scheduler: cancelling task for controller: %s\n", request.Controller.GetName())
		close(cancelChan)
	}
}

func (s *Scheduler) Run(stopChan chan struct{}) {
	for {
		select {
		case task := <-s.taskChan:
			//log.Printf("Scheduler: executing controller: %s\n", task.controller.GetName())
			reEnqueAfterSet := task.controller.Act(task.config)
			for _, reEnqueAfter := range reEnqueAfterSet {
				cancelItemChan := make(chan struct{})
				s.queue[reEnqueAfter.Controller.GetName()] = cancelItemChan
				go s.taskInQueue(reEnqueAfter, cancelItemChan)
			}
		case cancelRequest := <-s.cancelChan:
			cancelItemChan, ok := s.queue[cancelRequest.controller]
			if ok {
				s.cancelInner(cancelItemChan)
			} else {
				log.Printf("Scheduler: no task found to cancel for controller: %s\n", cancelRequest.controller)
			}
			delete(s.queue, cancelRequest.controller)
			close(cancelRequest.result)
		case <-stopChan:
			log.Println("Scheduler: shutting down. Canceling all tasks")

			for _, cancelItemChan := range s.queue {
				s.cancelInner(cancelItemChan)
			}
			return
		}
	}
}
