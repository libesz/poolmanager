package configstore

import (
	"io/ioutil"
	"log"

	"github.com/libesz/poolmanager/pkg/controller"
	"gopkg.in/yaml.v2"
)

func NewConfigStore(controllers []controller.Controller, hook ConfigStoreHook, backend ConfigStoreBackend) *ConfigStore {
	all, err := backend.Load()
	if err != nil {
		all = make(configSet)
		log.Println("ConfigStore: Could not load configuration, starting with defaults. Error:", err.Error())
		for _, controller := range controllers {
			all[controller.GetName()] = controller.GetDefaultConfig()
		}
		if err = backend.Save(all); err != nil {
			log.Println("ConfigStore: Could not save configuration, error:", err.Error())
		}
	}
	log.Printf("ConfigStore: loaded configuration: %+v", all)

	result := &ConfigStore{
		all:                    all,
		setChan:                make(chan configStoreSetArgs),
		getChan:                make(chan configStoreGetArgs),
		getControllerNamesChan: make(chan chan []string),
		getPropertiesChan:      make(chan configStoreGetPropertiesArgs),
		hook:                   hook,
		backend:                backend,
	}
	hook.SetConfigStore(result)
	return result
}

type configStoreFileBackend struct {
	path string
}

func NewConfigStoreFileBackend(path string) ConfigStoreBackend {
	return &configStoreFileBackend{
		path: path,
	}
}

func (f *configStoreFileBackend) Load() (configSet, error) {
	data := configSet{}

	rawContent, err := ioutil.ReadFile(f.path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(rawContent, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (f *configStoreFileBackend) Save(data configSet) error {
	rawContent, err := yaml.Marshal(&data)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(f.path, rawContent, 0600); err != nil {
		return err
	}
	return nil
}

func (s *ConfigStore) Set(controller string, config controller.Config, enqueue bool) error {
	resultChan := make(chan error)
	s.setChan <- configStoreSetArgs{controller: controller, config: config, resultChan: resultChan, enqueue: enqueue}
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

func (s *ConfigStore) GetControllerNames() []string {
	resultChan := make(chan []string)
	s.getControllerNamesChan <- resultChan
	return <-resultChan
}

func (s *ConfigStore) Run(stopChan chan struct{}) {
	for {
		select {
		case <-stopChan:
			return
		case getRequest := <-s.getChan:
			item, existing := s.all[getRequest.name]
			result := controller.EmptyConfig()
			if !existing {
				getRequest.resultChan <- result
				break
			}
			result = controller.CopyConfig(item)
			getRequest.resultChan <- result
		case getPropertiesRequest := <-s.getPropertiesChan:
			getPropertiesRequest.resultChan <- s.hook.GetConfigProperties(getPropertiesRequest.name)
		case setRequest := <-s.setChan:
			if orig, ok := s.all[setRequest.controller]; ok {
				if controller.IsEqualConfig(orig, setRequest.config) {
					setRequest.resultChan <- nil
					break
				}
			}
			err := s.hook.ConfigUpdated(setRequest.controller, setRequest.config, setRequest.enqueue)
			if err == nil {
				s.all[setRequest.controller] = setRequest.config
				s.backend.Save(s.all)
			}
			setRequest.resultChan <- err
		case getControllerNamesResponseChan := <-s.getControllerNamesChan:
			var result []string
			for controllerName := range s.all {
				result = append(result, controllerName)
			}
			getControllerNamesResponseChan <- result
		}
	}
}
