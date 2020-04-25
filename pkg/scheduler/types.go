package scheduler

import (
	"github.com/libesz/poolmanager/pkg/controller"
)

type Scheduler struct {
	taskChan chan controller.Controller
}
