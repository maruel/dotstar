// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package host

import (
	"errors"
	"fmt"
)

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
	Down         Pull = 1 // Apply pull-down
	Up           Pull = 2 // Apply pull-up
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
	// Function returns a user readable string representation of what the pin is
	// configured to do. Common case is In and Out but it can be bus specific pin
	// name.
	Function() string
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

// invalidPinErr is returned when trying to use INVALID.
var invalidPinErr = errors.New("invalid pin")

// INVALID implements PinIO for compability but fails on all access.
var INVALID invalidPin

// invalidPin implements PinIO for compability but fails on all access.
type invalidPin struct {
}

func (invalidPin) Number() int {
	return -1
}

func (invalidPin) String() string {
	return "INVALID"
}

func (invalidPin) Function() string {
	return "INVALID"
}

func (invalidPin) In(Pull) error {
	return invalidPinErr
}

func (invalidPin) Read() Level {
	return Low
}

func (invalidPin) Edges() (chan Level, error) {
	return nil, invalidPinErr
}

func (invalidPin) Out() error {
	return invalidPinErr
}

func (invalidPin) Set(Level) {
}
