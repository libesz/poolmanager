package configstore

import "github.com/libesz/poolmanager/pkg/controller"

type configSet map[string]controller.Config

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

type ConfigStoreBackend interface {
	Save(configSet) error
	Load() (configSet, error)
}

type ConfigStore struct {
	all                    configSet
	setChan                chan configStoreSetArgs
	getChan                chan configStoreGetArgs
	getControllerNamesChan chan chan []string
	getPropertiesChan      chan configStoreGetPropertiesArgs
	hook                   ConfigStoreHook
	backend                ConfigStoreBackend
}
