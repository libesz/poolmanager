package main

import (
	"sync"
	"time"

	"github.com/cloudflare/cfssl/log"
	"github.com/libesz/poolmanager/pkg/controller"
	"github.com/libesz/poolmanager/pkg/io"
	"github.com/libesz/poolmanager/pkg/scheduler"
	"github.com/libesz/poolmanager/pkg/webui"
)

func main() {
	pumpControllerConfig := controller.Config{"desired runtime per day": 1}
	pumpOutput := io.DummyOutput{Name: "pumpOutput"}
	timer := io.NewTimerOutput("pumpTimerOutput", &pumpOutput, time.Now)
	pumpOrOutputMembers := io.NewOrOutput(&timer, 2)
	pumpController := controller.NewPoolPumpController(&timer, &pumpOrOutputMembers[0])

	tempControllerConfig := controller.Config{"enabled": true, "desired temperature": 28.0, "start hour": 12, "end hour": 16}
	tempSensor := io.DummyTempSensor{Temperature: 26}
	heaterOutput := &io.DummyOutput{Name: "heater1"}
	tempController := controller.NewPoolTempController(0.5, &tempSensor, heaterOutput, &pumpOrOutputMembers[1], time.Now)

	stopChan := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	c := scheduler.NewConfigStore()
	go func() {
		c.Run(stopChan)
		wg.Done()
	}()

	s := scheduler.New(&c)
	wg.Add(1)
	go func() {
		s.Run(stopChan)
		wg.Done()
	}()
	if err := s.AddController(&tempController, tempControllerConfig); err != nil {
		log.Fatalf("Failed to add tempController: %s\n", err.Error())
	}
	if err := s.AddController(&pumpController, pumpControllerConfig); err != nil {
		log.Fatalf("Failed to add pumpController: %s\n", err.Error())
	}
	webui.Run(&c)
	wg.Wait()
}
