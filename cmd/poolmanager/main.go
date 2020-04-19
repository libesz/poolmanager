package main

import (
	"time"

	"github.com/libesz/poolmanager/pkg/controller"
	"github.com/libesz/poolmanager/pkg/io"
)

func hour9() time.Time {
	return time.Date(2020, 04, 15, 9, 0, 0, 0, time.Local)
}

func hour12() time.Time {
	return time.Date(2020, 04, 15, 12, 0, 0, 0, time.Local)
}

func hour14() time.Time {
	return time.Date(2020, 04, 15, 14, 0, 0, 0, time.Local)
}

func main() {
	timer := io.TimedGPIO{
		Name: "pump",
		Now:  hour9,
	}
	poolPumpControllerConfig := controller.Config{"desired runtime per day": 1}
	poolPumpController := controller.NewPoolPumpController(&timer)
	poolPumpController.Act(poolPumpControllerConfig)
	timer.Now = hour12
	poolPumpController.Act(poolPumpControllerConfig)

	poolTempControllerConfig := controller.Config{"desired temperature": 28, "start hour": 10, "end hour": 13}
	tempSensor := io.DummyTempSensor{Temperature: 26}
	heaterOutputs := []io.Output{
		&io.DummyOutput{Name: "heater1"},
		&io.DummyOutput{Name: "heater2"},
	}
	poolTempController := controller.PoolTempController{
		HeaterFactor:  0.5,
		HeaterOutputs: heaterOutputs,
		TempSensor:    &tempSensor,
		Now:           time.Now,
	}
	poolTempController.Act(poolTempControllerConfig)
}
