package scheduler

import (
	"github.com/libesz/poolmanager/pkg/configstore"
	"github.com/libesz/poolmanager/pkg/controller"
)

type Scheduler struct {
	taskChan    chan controller.Controller
	configStore *configstore.ConfigStore
}
