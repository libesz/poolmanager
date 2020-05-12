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
	pumpOutput := io.DummyOutput{Name_: "Pump state"}
	timer := io.NewTimerOutput("Pump runtime hours today", &pumpOutput, time.Now)
	pumpOrOutputMembers := io.NewOrOutput(&timer, 2)
	pumpController := controller.NewPoolPumpController(&timer, &pumpOrOutputMembers[0], time.Now)
	pumpControllerConfig := pumpController.GetDefaultConfig()

	tempSensor := io.DummyTempSensor{Temperature: 26}
	heaterOutput := &io.DummyOutput{Name_: "Heater"}
	tempController := controller.NewPoolTempController(0.5, &tempSensor, heaterOutput, &pumpOrOutputMembers[1], time.Now)
	tempControllerConfig := tempController.GetDefaultConfig()

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

	if err := c.Set(tempController.GetName(), tempControllerConfig, true); err != nil {
		log.Fatalf("Failed to set initial config for PoolTempController: %s\n", err.Error())
	}
	if err := c.Set(pumpController.GetName(), pumpControllerConfig, true); err != nil {
		log.Fatalf("Failed to set initial config for PoolPumpController: %s\n", err.Error())
	}

	wg.Add(1)
	w := webui.New(c, []io.Input{&tempSensor, &timer}, []io.Output{&pumpOutput, heaterOutput})
	go func() {
		w.Run(stopChan)
		wg.Done()
	}()

	wg.Wait()
}
