package scheduler

import (
	"time"

	"github.com/libesz/poolmanager/pkg/controller"
)

type Instance struct {
	controllers []controller.Controller
}

func New() Instance {
	return Instance{}
}

func (i *Instance) AddController(c controller.Controller) {
	i.controllers = append(i.controllers, c)
}

func (i *Instance) Run(pollTime time.Duration, stopChan chan struct{}) {
	ticker := time.NewTicker(pollTime).C
	for {
		select {
		case <-ticker:
			for _, c := range i.controllers {
				c.Act()
			}
		}
	}
}
