package main

import (
	"sync"
	"time"

	"github.com/libesz/poolmanager/pkg/controller"
	"github.com/libesz/poolmanager/pkg/io"
	"github.com/libesz/poolmanager/pkg/scheduler"
)

func main() {
	config := controller.Config{"desired runtime per day": 1, "desired temperature": 28, "start hour": 22, "end hour": 23}
	pumpOutput := io.DummyOutput{Name: "pumpOutput"}
	pumpOrOutputMembers := io.NewOrOutput(&pumpOutput, 2)

	timer := io.NewTimerOutput("pumpTimerOutput", &pumpOrOutputMembers[0], time.Now)
	poolPumpController := controller.NewPoolPumpController(&timer)

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
