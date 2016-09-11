// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package host

import "fmt"

// Level is the level of the pin: Low or High.
type Level bool

const (
	// Low represents 0v.
	Low Level = false
	// High represents Vin, generally 3.3v or 5v.
	High Level = true
)

func (l Level) String() string {
	if l == Low {
		return "Low"
	}
	return "High"
}

// Pull specifies the internal pull-up or pull-down for a pin set as input.
//
// The pull resistor stays set even after the processor shuts down. It is not
// possible to 'read back' what value was specified for each pin.
type Pull uint8

const (
	Float        Pull = 0 // Let the input float
	Down         Pull = 1 // Apply pull-down; for a bcm283x, the resistor is 50KOhm~60kOhm
	Up           Pull = 2 // Apply pull-up; for a bcm283x, the resistor is 50kOhm~65kOhm
	PullNoChange Pull = 3 // Do not change the previous pull resistor setting
)

const pullName = "FloatDownUpPullNoChange"

var pullIndex = [...]uint8{0, 5, 9, 11, 23}

func (i Pull) String() string {
	if i >= Pull(len(pullIndex)-1) {
		return fmt.Sprintf("Pull(%d)", i)
	}
	return pullName[pullIndex[i]:pullIndex[i+1]]
}

// Pin is a generic GPIO pin. Users are expected to use one of PinIn, PinOut or
// PinIO.
type Pin interface {
	fmt.Stringer
	// Number returns the logical pin number or a negative number if the pin is
	// not a GPIO, e.g. GROUND, V3_3, etc.
	Number() int
}

// PinIn is an input GPIO pin.
type PinIn interface {
	Pin

	// In setups a pin as an input.
	In(pull Pull) error
	// Read return the current pin level.
	//
	// Behavior is undefined if In() wasn't used before.
	Read() Level
	// Edges returns a channel that sends level changes.
	//
	// It is important to stop the querying loop by sending a Low to the channel
	// to stop it. The channel will then immediately be closed.
	//
	// If interrupt based edge detection is not supported, it will be emulated
	// via a query loop.
	//
	// Behavior is undefined if In() wasn't used before.
	Edges() (chan Level, error)
}

// PinOut is an output GPIO pin.
type PinOut interface {
	Pin

	// Out sets a pin as output. The caller should immediately call Set() after.
	Out() error
	// Set sets a pin already set for output as High or Low.
	//
	// Behavior is undefined if Out() wasn't used before.
	Set(l Level)
}

// PinIO is a GPIO pin that supports both input and output.
type PinIO interface {
	// PinIn and PinOut have to be duplicated here because Go doesn't allow
	// diamond shaped interfaces.
	Pin

	// In setups a pin as an input.
	In(pull Pull) error
	// Read return the current pin level.
	//
	// Behavior is undefined if In() wasn't used before.
	Read() Level
	// Edges returns a channel that sends level changes.
	//
	// It is important to stop the querying loop by sending a Low to the channel
	// to stop it. The channel will then immediately be closed.
	//
	// If interrupt based edge detection is not supported, it will be emulated
	// via a query loop.
	//
	// Behavior is undefined if In() wasn't used before.
	Edges() (chan Level, error)

	// Out sets a pin as output. The caller should immediately call Set() after.
	Out() error
	// Set sets a pin already set for output as High or Low.
	//
	// Behavior is undefined if Out() wasn't used before.
	Set(l Level)
}

// AllPins lists all the GPIO pins available on this host.
//
// This gets populated automatically on startup by the relevant child module,
// if running on a relevant host.
//
// Pins must be in order of their number.
//
// This list excludes non-GPIO pins like GROUND, V3_3, etc.
var AllPins []PinIO

// GetPinByName returns a GPIO pin from its name.
//
// This excludes non-GPIO pins like GROUND, V3_3, etc.
//
// Returns nil in case of failure.
func GetPinByName(name string) PinIO {
	// TODO(maruel): Create a map on first use?
	for _, p := range AllPins {
		if p.String() == name {
			return p
		}
	}
	return nil
}

// GetPinByNumber returns a GPIO pin from its number.
//
// This excludes non-GPIO pins like GROUND, V3_3, etc.
//
// Returns nil in case of failure.
func GetPinByNumber(number int) PinIO {
	if number >= 0 && number < len(AllPins) {
		return AllPins[number]
	}
	return nil
}