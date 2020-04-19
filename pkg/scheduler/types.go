package scheduler

import (
	"time"

	"github.com/libesz/poolmanager/pkg/controller"
)

type Interface interface {
	AddController(controller.Controller)
	Run(time.Duration, chan struct{})
}
