package scheduler

import "github.com/libesz/poolmanager/pkg/controller"

func NewConfigStore() ConfigStore {
	return ConfigStore{
		setChan:           make(chan configStoreSetArgs),
		getChan:           make(chan configStoreGetArgs),
		getKeysChan:       make(chan chan []string),
		getPropertiesChan: make(chan configStoreGetPropertiesArgs),
	}
}

func (s *ConfigStore) Set(controller controller.Controller, config controller.Config) error {
	resultChan := make(chan error)
	s.setChan <- configStoreSetArgs{controller: controller, config: config, resultChan: resultChan}
	return <-resultChan
}

func (s *ConfigStore) Get(name string) controller.Config {
	resultChan := make(chan controller.Config)
	s.getChan <- configStoreGetArgs{name: name, resultChan: resultChan}
	return <-resultChan
}

func (s *ConfigStore) GetProperties(name string) controller.ConfigProperties {
	resultChan := make(chan controller.ConfigProperties)
	s.getPropertiesChan <- configStoreGetPropertiesArgs{name: name, resultChan: resultChan}
	return <-resultChan
}

func (s *ConfigStore) GetKeys() []string {
	resultChan := make(chan []string)
	s.getKeysChan <- resultChan
	return <-resultChan
}

func (s *ConfigStore) Run(stopChan chan struct{}) {
	all := make(configSet)
	for {
		select {
		case <-stopChan:
			return
		case getRequest := <-s.getChan:
			item, existing := all[getRequest.name]
			if !existing {
				item = &configSetItem{}
			}
			getRequest.resultChan <- item.config
		case getPropertiesRequest := <-s.getPropertiesChan:
			item, existing := all[getPropertiesRequest.name]
			if !existing {
				getPropertiesRequest.resultChan <- controller.ConfigProperties{}
			}
			getPropertiesRequest.resultChan <- item.controller.GetConfig()
		case setRequest := <-s.setChan:
			err := setRequest.controller.ValidateConfig(setRequest.config)
			if err == nil {
				item := configSetItem{controller: setRequest.controller, config: setRequest.config}
				item.config = setRequest.config
				all[setRequest.controller.GetName()] = &item
			}
			setRequest.resultChan <- err
		case getKeysResponseChan := <-s.getKeysChan:
			var result []string
			for key := range all {
				result = append(result, key)
			}
			getKeysResponseChan <- result
		}
	}
}
