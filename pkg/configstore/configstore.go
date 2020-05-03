package configstore

import "github.com/libesz/poolmanager/pkg/controller"

func New(hook ConfigStoreHook) ConfigStore {
	return ConfigStore{
		setChan:           make(chan configStoreSetArgs),
		getChan:           make(chan configStoreGetArgs),
		getKeysChan:       make(chan chan []string),
		getPropertiesChan: make(chan configStoreGetPropertiesArgs),
		hook:              hook,
	}
}

func (s *ConfigStore) Set(controller string, config controller.Config) error {
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
			result := controller.EmptyConfig()
			if !existing {
				getRequest.resultChan <- result
				break
			}
			result = controller.CopyConfig(item.config)
			getRequest.resultChan <- result
		case getPropertiesRequest := <-s.getPropertiesChan:
			getPropertiesRequest.resultChan <- s.hook.GetConfigProperties(getPropertiesRequest.name)
		case setRequest := <-s.setChan:
			err := s.hook.ConfigUpdated(setRequest.controller, setRequest.config)
			if err == nil {
				item := configSetItem{controller: setRequest.controller, config: setRequest.config}
				item.config = setRequest.config
				all[setRequest.controller] = &item
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
