// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package sysfs

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/maruel/dlibox/go/pio"
	"github.com/maruel/dlibox/go/pio/conn/gpio"
)

// LEDs is all the leds discovered on this host via sysfs.
//
// Depending on the user context, the LEDs may be read-only or writeable.
var LEDs []*LED

// LEDByName returns a *LED for the LED name, if any.
//
// For all practical purpose, a LED is considered an output-only gpio.PinOut.
func LEDByName(name string) (*LED, error) {
	// TODO(maruel): Use a bisect or a map. For now we don't expect more than a
	// handful of LEDs so it doesn't matter.
	for _, led := range LEDs {
		if led.name == name {
			if err := led.open(); err != nil {
				return nil, err
			}
			return led, nil
		}
	}
	return nil, errors.New("invalid LED name")
}

// LED represents one LED on the system.
type LED struct {
	number int
	name   string
	root   string

	lock        sync.Mutex
	fBrightness *os.File // handle to /sys/class/gpio/gpio*/direction; never closed
}

func (l *LED) String() string {
	return l.name
}

// Number implements pins.Pin.
func (l *LED) Number() int {
	return l.number
}

// Function implements pins.Pin.
func (l *LED) Function() string {
	if l.Read() {
		return "LED/On"
	}
	return "LED/Off"
}

// In implements gpio.PinIn.
func (l *LED) In(pull gpio.Pull, edge gpio.Edge) error {
	if pull != gpio.Float && pull != gpio.PullNoChange {
		return errors.New("pull is not supported on LED")
	}
	if edge != gpio.None {
		return errors.New("edge is not supported on LED")
	}
	return nil
}

// Read implements gpio.PinIn.
func (l *LED) Read() gpio.Level {
	err := l.open()
	if err != nil {
		return gpio.Low
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	if _, err := l.fBrightness.Seek(0, 0); err != nil {
		return gpio.Low
	}
	var b [4]byte
	if _, err := l.fBrightness.Read(b[:]); err != nil {
		return gpio.Low
	}
	if b[0] != '0' {
		return gpio.High
	}
	return gpio.Low
}

// WaitForEdge implements gpio.PinIn.
func (l *LED) WaitForEdge(timeout time.Duration) bool {
	return false
}

// Pull implements gpio.PinIn.
func (l *LED) Pull() gpio.Pull {
	return gpio.PullNoChange
}

// Out implements gpio.PinOut.
func (l *LED) Out(level gpio.Level) error {
	err := l.open()
	if err != nil {
		return err
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	if _, err = l.fBrightness.Seek(0, 0); err != nil {
		return err
	}
	if level {
		_, err = l.fBrightness.Write([]byte("255"))
	} else {
		_, err = l.fBrightness.Write([]byte("0"))
	}
	return err
}

// PWM implements gpio.PinOut.
func (l *LED) PWM(duty int) error {
	return errors.New("pwm is not supported")
}

//

func (l *LED) open() error {
	l.lock.Lock()
	defer l.lock.Unlock()
	// trigger, max_brightness.
	var err error
	if l.fBrightness == nil {
		p := l.root + "brightness"
		if l.fBrightness, err = os.OpenFile(p, os.O_RDWR, 0600); err != nil {
			// Retry with read-only. This is the default setting.
			l.fBrightness, err = os.OpenFile(p, os.O_RDONLY, 0600)
		}
	}
	return err
}

// driverLED implements pio.Driver.
type driverLED struct {
}

func (d *driverLED) String() string {
	return "sysfs-led"
}

func (d *driverLED) Type() pio.Type {
	return pio.Pins
}

func (d *driverLED) Prerequisites() []string {
	return nil
}

// Init initializes LEDs sysfs handling code.
//
// Uses led sysfs as described* at
// https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-class-led
//
// * for the most minimalistic meaning of 'described'.
func (d *driverLED) Init() (bool, error) {
	items, err := filepath.Glob("/sys/class/leds/*")
	if err != nil {
		return true, err
	}
	if len(items) == 0 {
		return false, errors.New("no LED found")
	}
	// This make the LEDs in deterministic order.
	sort.Strings(items)
	for i, item := range items {
		LEDs = append(LEDs, &LED{
			number: i,
			name:   filepath.Base(item),
			root:   item + "/",
		})
	}
	return true, nil
}

func init() {
	if isLinux {
		pio.MustRegister(&driverLED{})
	}
}

var _ gpio.PinIn = &LED{}
var _ gpio.PinOut = &LED{}
var _ gpio.PinIO = &LED{}
