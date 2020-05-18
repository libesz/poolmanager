package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/libesz/poolmanager/pkg/configstore"
	"github.com/libesz/poolmanager/pkg/controller"
	"github.com/libesz/poolmanager/pkg/io"
	"github.com/libesz/poolmanager/pkg/scheduler"
	"github.com/libesz/poolmanager/pkg/webui"
)

func cleanup(cleanTheseUp []io.Haltable) {
	log.Println("Main: cleaning up...")
	for _, haltable := range cleanTheseUp {
		haltable.Halt()
	}
}

type StaticConfig struct {
	ListenOn      string
	Password      string
	PumpGPIO1     string
	PumpGPIO2     string
	HeaterGPIO    string
	HeaterOffDuty string
	HeaterOnDuty  string
}

func main() {
	sigs := make(chan os.Signal, 1)
	signalReceived := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println()
		log.Println("Main: signal received:", sig)
		signalReceived <- true
	}()

	var staticConfig StaticConfig
	if err := envconfig.Process("poolmanager", &staticConfig); err != nil {
		log.Fatal(err.Error())
	}

	var cleanTheseUp []io.Haltable
	defer func() {
		cleanup(cleanTheseUp)
	}()

	//pumpOutput1 := io.DummyOutput{Name_: "Pump1"}
	//pumpOutput2 := io.DummyOutput{Name_: "Pump2"}
	pumpOutput1 := io.NewGPIOOutput("Pump1", staticConfig.PumpGPIO1, true)
	cleanTheseUp = append(cleanTheseUp, pumpOutput1)
	pumpOutput2 := io.NewGPIOOutput("Pump2", staticConfig.PumpGPIO2, true)
	cleanTheseUp = append(cleanTheseUp, pumpOutput2)
	pumpOutput := io.NewOutputDistributor("Pump", []io.Output{pumpOutput1, pumpOutput2})
	timer := io.NewTimerOutput("Pump runtime hours today", pumpOutput, time.Now)
	pumpOrOutputMembers := io.NewOrOutput("Pump", &timer, 2)
	pumpController := controller.NewPoolPumpController(&timer, &pumpOrOutputMembers[0], time.Now)
	pumpControllerConfig := pumpController.GetDefaultConfig()

	tempSensor := io.DummyTempSensor{Temperature: 26}
	//heaterOutput := &io.DummyOutput{Name_: "Heater"}
	heaterOutput := io.NewGPIOOutput("Heater", staticConfig.HeaterGPIO, true)
	cleanTheseUp = append(cleanTheseUp, heaterOutput)
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
		log.Fatalf("Main: failed to set initial config for PoolTempController: %s\n", err.Error())
	}
	if err := c.Set(pumpController.GetName(), pumpControllerConfig, true); err != nil {
		log.Fatalf("Main: failed to set initial config for PoolPumpController: %s\n", err.Error())
	}

	wg.Add(1)
	w := webui.New(staticConfig.ListenOn, staticConfig.Password, c, []io.Input{&tempSensor, &timer}, []io.Output{pumpOutput, heaterOutput})
	go func() {
		w.Run(stopChan)
		wg.Done()
	}()

	<-signalReceived
	close(stopChan)
	wg.Wait()
	log.Println("Main: Exiting")
}
