// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package sysfs

// MakeSPI opens a *SPI via its sysfs interface as described at
// https://www.kernel.org/doc/Documentation/spi/spidev.
//
// `speed` must be specified and should be in the high Khz or low Mhz range,
// it's a good idea to start at 4000000 (4Mhz) and go upward as long as the
// signal is good.
//
// Default configuration is Mode3 and 8 bits.
func MakeSPI(bus, chipSelect int, speed int64) (*SPI, error) {
	return makeSPI(bus, chipSelect, speed)
}
