package main

import (
	"sync"
	"time"

	"github.com/libesz/poolmanager/pkg/controller"
	"github.com/libesz/poolmanager/pkg/io"
	"github.com/libesz/poolmanager/pkg/scheduler"
)

func main() {
	pumpControllerConfig := controller.Config{"desired runtime per day": 1}
	pumpOutput := io.DummyOutput{Name: "pumpOutput"}
	timer := io.NewTimerOutput("pumpTimerOutput", &pumpOutput, time.Now)
	pumpOrOutputMembers := io.NewOrOutput(&timer, 2)
	pumpController := controller.NewPoolPumpController(&timer, &pumpOrOutputMembers[0])

	tempControllerConfig := controller.Config{"desired temperature": 28, "start hour": 12, "end hour": 16}
	tempSensor := io.DummyTempSensor{Temperature: 26}
	heaterOutput := &io.DummyOutput{Name: "heater1"}
	tempController := controller.NewPoolTempController(0.5, &tempSensor, heaterOutput, &pumpOrOutputMembers[1], time.Now)

	s := scheduler.New()
	stopChan := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		s.Run(stopChan)
		wg.Done()
	}()
	s.AddController(&tempController, &tempControllerConfig)
	s.AddController(&pumpController, &pumpControllerConfig)
	wg.Wait()
}
