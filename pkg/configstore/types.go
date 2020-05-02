package configstore

import "github.com/libesz/poolmanager/pkg/controller"

type configSetItem struct {
	controller controller.Controller
	config     controller.Config
}

type configSet map[string]*configSetItem

type configStoreSetArgs struct {
	controller controller.Controller
	config     controller.Config
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

type ConfigStore struct {
	setChan           chan configStoreSetArgs
	getChan           chan configStoreGetArgs
	getKeysChan       chan chan []string
	getPropertiesChan chan configStoreGetPropertiesArgs
}
