package scheduler

import (
	"log"
	"time"

	"github.com/libesz/poolmanager/pkg/controller"
)

func New() Scheduler {
	return Scheduler{taskChan: make(chan controller.Controller), configSet: make(ConfigSet)}
}

func (s *Scheduler) AddController(c controller.Controller, config *controller.Config) {
	log.Printf("Scheduler: added controller: %s\n", c.GetName())
	s.configSet[c.GetName()] = config
	s.enqueue(c)
}

func (s *Scheduler) enqueue(c controller.Controller) {
	s.taskChan <- c
}

func (s *Scheduler) Run(stopChan chan struct{}) {
	for {
		select {
		case c := <-s.taskChan:
			log.Printf("Scheduler: executing controller: %s\n", c.GetName())
			config := &controller.Config{}
			if configFromSet, ok := s.configSet[c.GetName()]; ok {
				config = configFromSet
			}
			reEnqueAfterSet := c.Act(*config)
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
