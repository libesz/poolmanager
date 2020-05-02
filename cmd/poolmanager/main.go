package main

import (
	"log"
	"sync"
	"time"

	"github.com/libesz/poolmanager/pkg/configstore"
	"github.com/libesz/poolmanager/pkg/controller"
	"github.com/libesz/poolmanager/pkg/io"
	"github.com/libesz/poolmanager/pkg/scheduler"
	"github.com/libesz/poolmanager/pkg/webui"
)

func main() {
	pumpControllerConfig := controller.Config{"desired runtime per day": 1.0}
	pumpOutput := io.DummyOutput{Name: "pumpOutput"}
	timer := io.NewTimerOutput("pumpTimerOutput", &pumpOutput, time.Now)
	pumpOrOutputMembers := io.NewOrOutput(&timer, 2)
	pumpController := controller.NewPoolPumpController(&timer, &pumpOrOutputMembers[0])

	tempControllerConfig := controller.Config{"enabled": true, "desired temperature": 28.0, "start hour": 12.0, "end hour": 16.0}
	tempSensor := io.DummyTempSensor{Temperature: 26}
	heaterOutput := &io.DummyOutput{Name: "heater1"}
	tempController := controller.NewPoolTempController(0.5, &tempSensor, heaterOutput, &pumpOrOutputMembers[1], time.Now)

	stopChan := make(chan struct{})
	wg := sync.WaitGroup{}

	s := scheduler.New()
	wg.Add(1)
	go func() {
		s.Run(stopChan)
		wg.Done()
	}()

	wg.Add(1)
	c := configstore.New(&s)
	go func() {
		c.Run(stopChan)
		wg.Done()
	}()

	s.AddController(&tempController)
	s.AddController(&pumpController)

	if err := c.Set(tempController.GetName(), tempControllerConfig); err != nil {
		log.Printf("Failed to set initial config for PoolTempController: %s\n", err.Error())
	}
	if err := c.Set(pumpController.GetName(), pumpControllerConfig); err != nil {
		log.Printf("Failed to set initial config for PoolPumpController: %s\n", err.Error())
	}

	webui.Run(&c)
	wg.Wait()
}
