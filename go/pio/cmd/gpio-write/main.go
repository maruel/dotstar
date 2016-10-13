// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// gpio-write sets a GPIO pin to low or high.
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/maruel/dlibox/go/pio/conn/gpio"
	"github.com/maruel/dlibox/go/pio/host"
)

func mainImpl() error {
	if len(os.Args) != 3 {
		return errors.New("specify GPIO pin to write to and its level (0 or 1)")
	}
	pin, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return err
	}

	l := gpio.Low
	switch os.Args[2] {
	case "0":
	case "1":
		l = gpio.High
	default:
		return errors.New("specify level as 0 or 1")
	}

	if _, err := host.Init(); err != nil {
		return err
	}

	p := gpio.ByNumber(pin)
	if p == nil {
		return errors.New("invalid GPIO pin number")
	}

	return p.Out(l)
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "gpio-write: %s.\n", err)
		os.Exit(1)
	}
}
