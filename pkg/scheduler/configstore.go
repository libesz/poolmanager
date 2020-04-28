package scheduler

import "github.com/libesz/poolmanager/pkg/controller"

func NewConfigStore() ConfigStore {
	return ConfigStore{
		setChan:     make(chan configStoreSetArgs),
		getChan:     make(chan configStoreGetArgs),
		getKeysChan: make(chan chan []string),
	}
}

func (s *ConfigStore) Set(name string, config controller.Config) {
	s.setChan <- configStoreSetArgs{name: name, config: config}
}

func (s *ConfigStore) Get(name string) controller.Config {
	resultChan := make(chan controller.Config)
	s.getChan <- configStoreGetArgs{name: name, resultChan: resultChan}
	return <-resultChan
}

func (s *ConfigStore) GetKeys() []string {
	resultChan := make(chan []string)
	s.getKeysChan <- resultChan
	return <-resultChan
}

func (s *ConfigStore) Run(stopChan chan struct{}) {
	all := configSet{}
	for {
		select {
		case <-stopChan:
			return
		case getRequest := <-s.getChan:
			getRequest.resultChan <- all[getRequest.name]
		case setRequest := <-s.setChan:
			all[setRequest.name] = setRequest.config
		case getKeysResponseChan := <-s.getKeysChan:
			var result []string
			for key := range all {
				result = append(result, key)
			}
			getKeysResponseChan <- result
		}
	}
}
