package scheduler

import (
	"github.com/libesz/poolmanager/pkg/controller"
)

type ConfigSet map[string]*controller.Config

type Scheduler struct {
	taskChan  chan controller.Controller
	configSet ConfigSet
}
