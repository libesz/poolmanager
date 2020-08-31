package scheduler

import (
	"sync"
	"testing"

	"github.com/libesz/poolmanager/pkg/controller"
)

type MockController struct {
	ActReturnThis chan []controller.EnqueueRequest
	ActStarted    chan struct{}
}

func (c *MockController) Act(config controller.Config) []controller.EnqueueRequest {
	var returnThis []controller.EnqueueRequest
	if c.ActReturnThis != nil {
		returnThis = <-c.ActReturnThis
	}
	if c.ActStarted != nil {
		c.ActStarted <- struct{}{}
	}

	return returnThis
}

func (c *MockController) GetName() string {
	return "mock"
}

func (c MockController) GetConfigProperties() controller.ConfigProperties {
	return controller.ConfigProperties{}
}

func (c MockController) ValidateConfig(controller.Config) error {
	return nil
}

func (c MockController) GetDefaultConfig() controller.Config {
	return controller.EmptyConfig()
}

func TestDefault(t *testing.T) {
	sch := New()
	stopChan := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		sch.Run(stopChan)
		wg.Done()
	}()
	mockController := MockController{ActStarted: make(chan struct{})}
	sch.AddController(&mockController)
	<-mockController.ActStarted
	close(stopChan)
	wg.Wait()
}

func TestReenqueue(t *testing.T) {
	sch := New()
	stopChan := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		sch.Run(stopChan)
		wg.Done()
	}()
	mockController := MockController{ActStarted: make(chan struct{}), ActReturnThis: make(chan []controller.EnqueueRequest)}
	sch.AddController(&mockController)
	mockController.ActReturnThis <- []controller.EnqueueRequest{{
		Controller: &mockController,
		After:      0,
	}}
	<-mockController.ActStarted
	mockController.ActReturnThis <- []controller.EnqueueRequest{}
	<-mockController.ActStarted
	close(stopChan)
	wg.Wait()
}

func TestConfigUpdated(t *testing.T) {
	sch := New()
	stopChan := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		sch.Run(stopChan)
		wg.Done()
	}()
	mockController := MockController{ActStarted: make(chan struct{}), ActReturnThis: make(chan []controller.EnqueueRequest)}
	sch.AddController(&mockController)
	mockController.ActReturnThis <- []controller.EnqueueRequest{{
		Controller: &mockController,
		After:      0,
	}}
	<-mockController.ActStarted
	sch.ConfigUpdated(mockController.GetName(), controller.EmptyConfig(), true)
	mockController.ActReturnThis <- []controller.EnqueueRequest{}
	<-mockController.ActStarted
	close(stopChan)
	wg.Wait()
}

func TestConfigUpdatedForNonExistingController(t *testing.T) {
	sch := New()
	stopChan := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		sch.Run(stopChan)
		wg.Done()
	}()
	if sch.ConfigUpdated("does not exist", controller.EmptyConfig(), true) == nil {
		t.Fatal("Config change shall result in error for non-queued controller")
	}
	close(stopChan)
	wg.Wait()
}
