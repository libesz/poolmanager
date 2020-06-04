package io

import "math"

//InputError represents that the value is currently not available from the device
const InputError = math.MaxFloat32

//Input represents any measurable property as a floating point number
type Input interface {
	//Name tells the human readable name of the input
	Name() string
	//Type tells the human readable type / nature of the input (i.e. "Temperature")
	Type() string
	//Degree tells the human readable degree of the input (i.e. "°C")
	Degree() string
	//Value tells the actual up-to-date value of the measured indicator
	Value() float64
}

//Output represents any controlled device which receives a single ON of OFF signal
type Output interface {
	//Name tells the human readable name of the output
	Name() string
	//Set changes the device state. It returns a flag wether the change was effective (i.e. returns false if nothing changed)
	Set(bool) bool
	//Get tells the actual device state
	Get() bool
}

//Haltable is an interface which is expected to be implemented by outputs which need careful cleanup on system shutdown
type Haltable interface {
	//Halt shuts down the device as appropriate. At least it shall Set() the device into the idle state.
	Halt()
}
