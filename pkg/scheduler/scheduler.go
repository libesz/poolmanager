package scheduler

import (
	"log"
	"time"

	"github.com/libesz/poolmanager/pkg/configstore"
	"github.com/libesz/poolmanager/pkg/controller"
)

func New(configStore *configstore.ConfigStore) Scheduler {
	return Scheduler{taskChan: make(chan controller.Controller), configStore: configStore}
}

func (s *Scheduler) AddController(c controller.Controller, config controller.Config) error {
	log.Printf("Scheduler: added controller: %s\n", c.GetName())
	err := s.configStore.Set(c, config)
	if err != nil {
		return err
	}
	s.enqueue(c)
	return nil
}

func (s *Scheduler) enqueue(c controller.Controller) {
	s.taskChan <- c
}

func (s *Scheduler) Run(stopChan chan struct{}) {
	for {
		select {
		case c := <-s.taskChan:
			log.Printf("Scheduler: executing controller: %s\n", c.GetName())
			config := controller.Config{}
			if configFromSet := s.configStore.Get(c.GetName()); configFromSet != nil {
				config = configFromSet
			}
			reEnqueAfterSet := c.Act(config)
			for _, reEnqueAfter := range reEnqueAfterSet {
				go func(request controller.EnqueueRequest) {
					time.Sleep(request.After)
					s.enqueue(request.Controller)
				}(reEnqueAfter)
			}
		case <-stopChan:
			return
		}
	}
}
