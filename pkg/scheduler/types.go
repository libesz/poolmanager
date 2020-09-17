package scheduler

import (
	"github.com/libesz/poolmanager/pkg/configstore"
	"github.com/libesz/poolmanager/pkg/controller"
)

type schedulerTask struct {
	controller controller.Controller
	config     controller.Config
}

type cancelTask struct {
	controller string
	result     chan struct{}
}

type Scheduler struct {
	taskChan    chan schedulerTask
	cancelChan  chan cancelTask
	queue       map[string]chan chan struct{}
	configStore *configstore.ConfigStore
	controllers map[string]controller.Controller
}
