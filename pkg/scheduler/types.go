package scheduler

import (
	"github.com/libesz/poolmanager/pkg/controller"
)

type configSet map[string]controller.Config

type configStoreSetArgs struct {
	name   string
	config controller.Config
}

type configStoreGetArgs struct {
	name       string
	resultChan chan controller.Config
}

type ConfigStore struct {
	setChan chan configStoreSetArgs
	getChan chan configStoreGetArgs
}

type Scheduler struct {
	taskChan    chan controller.Controller
	configStore *ConfigStore
}
