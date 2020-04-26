package main

import (
	"sync"
	"time"

	"github.com/libesz/poolmanager/pkg/controller"
	"github.com/libesz/poolmanager/pkg/io"
	"github.com/libesz/poolmanager/pkg/scheduler"
)

func main() {
	config := controller.Config{"desired runtime per day": 1, "desired temperature": 28, "start hour": 12, "end hour": 15}
	pumpOutput := io.DummyOutput{Name: "pumpOutput"}
	timer := io.NewTimerOutput("pumpTimerOutput", &pumpOutput, time.Now)
	pumpOrOutputMembers := io.NewOrOutput(&timer, 2)
	poolPumpController := controller.NewPoolPumpController(&timer, &pumpOrOutputMembers[0])

	tempSensor := io.DummyTempSensor{Temperature: 26}
	heaterOutput := &io.DummyOutput{Name: "heater1"}
	poolTempController := controller.NewPoolTempController(0.5, &tempSensor, heaterOutput, &pumpOrOutputMembers[1], time.Now)

	s := scheduler.New()
	stopChan := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		s.Run(&config, stopChan)
		wg.Done()
	}()
	s.AddController(&poolTempController)
	s.AddController(&poolPumpController)
	wg.Wait()
}
