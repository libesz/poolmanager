package configstore

import "github.com/libesz/poolmanager/pkg/controller"

type configSetItem struct {
	controller string
	config     controller.Config
}

type configSet map[string]*configSetItem

type configStoreSetArgs struct {
	controller string
	config     controller.Config
	enqueue    bool
	resultChan chan error
}

type configStoreGetPropertiesArgs struct {
	name       string
	resultChan chan controller.ConfigProperties
}

type configStoreGetArgs struct {
	name       string
	resultChan chan controller.Config
}

type ConfigStoreHook interface {
	ConfigUpdated(controller string, config controller.Config, enqueue bool) error
	GetConfigProperties(controller string) controller.ConfigProperties
	SetConfigStore(configStore *ConfigStore)
}

type ConfigStore struct {
	setChan           chan configStoreSetArgs
	getChan           chan configStoreGetArgs
	getKeysChan       chan chan []string
	getPropertiesChan chan configStoreGetPropertiesArgs
	hook              ConfigStoreHook
}
