// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"time"

	"github.com/mattn/go-colorable"
)

type dotStar struct {
	w          io.WriteCloser
	b          []byte
	brightness int // 0-31
}

func (d *dotStar) Close() error {
	return d.w.Close()
}

func (d *dotStar) Write(pixels []color.NRGBA) error {
	// https://cpldcpu.files.wordpress.com/2014/08/apa-102c-super-led-specifications-2014-en.pdf
	numLights := len(pixels)
	// End frames are needed to be able to push enough SPI clock signals due to
	// internal half-delay of data signal from each individual LED. See
	// https://cpldcpu.wordpress.com/2014/11/30/understanding-the-apa102-superled/
	l := 4*(numLights+1) + numLights/2/8 + 1
	if len(d.b) < l {
		d.b = make([]byte, l)
	}
	// Start frame is all zeros. Just skip it.
	s := d.b[4:]
	brightness := byte(0xE0 + d.brightness)
	for i := range pixels {
		r, g, b, _ := pixels[i].RGBA()
		// BGR.
		s[4*i] = brightness
		s[4*i+1] = byte(b >> 8)
		s[4*i+2] = byte(g >> 8)
		s[4*i+3] = byte(r >> 8)
	}
	// End frames
	s = s[4*numLights:]
	for i := range s {
		s[i] = 0xFF
	}
	_, err := d.w.Write(d.b)
	return err
}

func (d *dotStar) MinDelay() time.Duration {
	// As per APA102-C spec, it's max refresh rate is 400hz.
	// https://en.wikipedia.org/wiki/Flicker_fusion_threshold is a recommended
	// reading.
	return time.Second / 400
}

// MakeDotStar returns a stripe that communicates over SPI.
func MakeDotStar() (Strip, error) {
	// The speed must be high, as there's 32 bits sent per LED, creating a
	// staggered effect. See
	// https://cpldcpu.wordpress.com/2014/11/30/understanding-the-apa102-superled/
	w, err := MakeSPI("", 20000000)
	if err != nil {
		return nil, err
	}
	return &dotStar{w: w, brightness: 31}, err
}

//

type screenStrip struct {
	w io.Writer
	b bytes.Buffer
}

func (s *screenStrip) Close() error {
	return nil
}

func (s *screenStrip) Write(pixels []color.NRGBA) error {
	// This code is designed to minimize the amount of memory allocated per call.
	s.b.Reset()
	_, _ = s.b.WriteString("\r\033[0m")
	lastI := -1
	for _, c := range pixels {
		newI := rgbToANSI(c)
		if newI != lastI {
			// Only send the ANSI code when the color changes.
			lastI = newI
			_, _ = fmt.Fprintf(&s.b, "\033[48;5;%dm", newI)
		}
		_, _ = s.b.WriteString(" ")
	}
	_, _ = s.b.WriteString("\033[0m ")
	_, err := s.b.WriteTo(s.w)
	return err
}

func (s *screenStrip) MinDelay() time.Duration {
	// Limit to 30hz, especially for ssh connections.
	return time.Second / 30
}

// MakeScreen returns a stripe that display at the console.
func MakeScreen() Strip {
	return &screenStrip{w: colorable.NewColorableStdout()}
}

//

// lookupANSIToRGB is a look up table for ANSI color codes. It is incorrect but
// for now #closeenough to test while waiting for the hardware to arrive from
// China.
var lookupANSIToRGB = [256]color.NRGBA{
	{0x00, 0x00, 0x00, 0xFF},
	{0x80, 0x00, 0x00, 0xFF},
	{0x00, 0x80, 0x00, 0xFF},
	{0x80, 0x80, 0x00, 0xFF},
	{0x00, 0x00, 0x80, 0xFF},
	{0x80, 0x00, 0x80, 0xFF},
	{0x00, 0x80, 0x80, 0xFF},
	{0xC0, 0xC0, 0xC0, 0xFF},
	{0x80, 0x80, 0x80, 0xFF},
	{0xFF, 0x00, 0x00, 0xFF},
	{0x00, 0xFF, 0x00, 0xFF},
	{0xFF, 0xFF, 0x00, 0xFF},
	{0x00, 0x00, 0xFF, 0xFF},
	{0xFF, 0x00, 0xFF, 0xFF},
	{0x00, 0xFF, 0xFF, 0xFF},
	{0xFF, 0xFF, 0xFF, 0xFF},
	{0x00, 0x00, 0x00, 0xFF},
	{0x00, 0x00, 0x5F, 0xFF},
	{0x00, 0x00, 0x87, 0xFF},
	{0x00, 0x00, 0xAF, 0xFF},
	{0x00, 0x00, 0xD7, 0xFF},
	{0x00, 0x00, 0xFF, 0xFF},
	{0x00, 0x5F, 0x00, 0xFF},
	{0x00, 0x5F, 0x5F, 0xFF},
	{0x00, 0x5F, 0x87, 0xFF},
	{0x00, 0x5F, 0xAF, 0xFF},
	{0x00, 0x5F, 0xD7, 0xFF},
	{0x00, 0x5F, 0xFF, 0xFF},
	{0x00, 0x87, 0x00, 0xFF},
	{0x00, 0x87, 0x5F, 0xFF},
	{0x00, 0x87, 0x87, 0xFF},
	{0x00, 0x87, 0xAF, 0xFF},
	{0x00, 0x87, 0xD7, 0xFF},
	{0x00, 0x87, 0xFF, 0xFF},
	{0x00, 0xAF, 0x00, 0xFF},
	{0x00, 0xAF, 0x5F, 0xFF},
	{0x00, 0xAF, 0x87, 0xFF},
	{0x00, 0xAF, 0xAF, 0xFF},
	{0x00, 0xAF, 0xD7, 0xFF},
	{0x00, 0xAF, 0xFF, 0xFF},
	{0x00, 0xD7, 0x00, 0xFF},
	{0x00, 0xD7, 0x5F, 0xFF},
	{0x00, 0xD7, 0x87, 0xFF},
	{0x00, 0xD7, 0xAF, 0xFF},
	{0x00, 0xD7, 0xD7, 0xFF},
	{0x00, 0xD7, 0xFF, 0xFF},
	{0x00, 0xFF, 0x00, 0xFF},
	{0x00, 0xFF, 0x5F, 0xFF},
	{0x00, 0xFF, 0x87, 0xFF},
	{0x00, 0xFF, 0xAF, 0xFF},
	{0x00, 0xFF, 0xD7, 0xFF},
	{0x00, 0xFF, 0xFF, 0xFF},
	{0x5F, 0x00, 0x00, 0xFF},
	{0x5F, 0x00, 0x5F, 0xFF},
	{0x5F, 0x00, 0x87, 0xFF},
	{0x5F, 0x00, 0xAF, 0xFF},
	{0x5F, 0x00, 0xD7, 0xFF},
	{0x5F, 0x00, 0xFF, 0xFF},
	{0x5F, 0x5F, 0x00, 0xFF},
	{0x5F, 0x5F, 0x5F, 0xFF},
	{0x5F, 0x5F, 0x87, 0xFF},
	{0x5F, 0x5F, 0xAF, 0xFF},
	{0x5F, 0x5F, 0xD7, 0xFF},
	{0x5F, 0x5F, 0xFF, 0xFF},
	{0x5F, 0x87, 0x00, 0xFF},
	{0x5F, 0x87, 0x5F, 0xFF},
	{0x5F, 0x87, 0x87, 0xFF},
	{0x5F, 0x87, 0xAF, 0xFF},
	{0x5F, 0x87, 0xD7, 0xFF},
	{0x5F, 0x87, 0xFF, 0xFF},
	{0x5F, 0xAF, 0x00, 0xFF},
	{0x5F, 0xAF, 0x5F, 0xFF},
	{0x5F, 0xAF, 0x87, 0xFF},
	{0x5F, 0xAF, 0xAF, 0xFF},
	{0x5F, 0xAF, 0xD7, 0xFF},
	{0x5F, 0xAF, 0xFF, 0xFF},
	{0x5F, 0xD7, 0x00, 0xFF},
	{0x5F, 0xD7, 0x5F, 0xFF},
	{0x5F, 0xD7, 0x87, 0xFF},
	{0x5F, 0xD7, 0xAF, 0xFF},
	{0x5F, 0xD7, 0xD7, 0xFF},
	{0x5F, 0xD7, 0xFF, 0xFF},
	{0x5F, 0xFF, 0x00, 0xFF},
	{0x5F, 0xFF, 0x5F, 0xFF},
	{0x5F, 0xFF, 0x87, 0xFF},
	{0x5F, 0xFF, 0xAF, 0xFF},
	{0x5F, 0xFF, 0xD7, 0xFF},
	{0x5F, 0xFF, 0xFF, 0xFF},
	{0x87, 0x00, 0x00, 0xFF},
	{0x87, 0x00, 0x5F, 0xFF},
	{0x87, 0x00, 0x87, 0xFF},
	{0x87, 0x00, 0xAF, 0xFF},
	{0x87, 0x00, 0xD7, 0xFF},
	{0x87, 0x00, 0xFF, 0xFF},
	{0x87, 0x5F, 0x00, 0xFF},
	{0x87, 0x5F, 0x5F, 0xFF},
	{0x87, 0x5F, 0x87, 0xFF},
	{0x87, 0x5F, 0xAF, 0xFF},
	{0x87, 0x5F, 0xD7, 0xFF},
	{0x87, 0x5F, 0xFF, 0xFF},
	{0x87, 0x87, 0x00, 0xFF},
	{0x87, 0x87, 0x5F, 0xFF},
	{0x87, 0x87, 0x87, 0xFF},
	{0x87, 0x87, 0xAF, 0xFF},
	{0x87, 0x87, 0xD7, 0xFF},
	{0x87, 0x87, 0xFF, 0xFF},
	{0x87, 0xAF, 0x00, 0xFF},
	{0x87, 0xAF, 0x5F, 0xFF},
	{0x87, 0xAF, 0x87, 0xFF},
	{0x87, 0xAF, 0xAF, 0xFF},
	{0x87, 0xAF, 0xD7, 0xFF},
	{0x87, 0xAF, 0xFF, 0xFF},
	{0x87, 0xD7, 0x00, 0xFF},
	{0x87, 0xD7, 0x5F, 0xFF},
	{0x87, 0xD7, 0x87, 0xFF},
	{0x87, 0xD7, 0xAF, 0xFF},
	{0x87, 0xD7, 0xD7, 0xFF},
	{0x87, 0xD7, 0xFF, 0xFF},
	{0x87, 0xFF, 0x00, 0xFF},
	{0x87, 0xFF, 0x5F, 0xFF},
	{0x87, 0xFF, 0x87, 0xFF},
	{0x87, 0xFF, 0xAF, 0xFF},
	{0x87, 0xFF, 0xD7, 0xFF},
	{0x87, 0xFF, 0xFF, 0xFF},
	{0xAF, 0x00, 0x00, 0xFF},
	{0xAF, 0x00, 0x5F, 0xFF},
	{0xAF, 0x00, 0x87, 0xFF},
	{0xAF, 0x00, 0xAF, 0xFF},
	{0xAF, 0x00, 0xD7, 0xFF},
	{0xAF, 0x00, 0xFF, 0xFF},
	{0xAF, 0x5F, 0x00, 0xFF},
	{0xAF, 0x5F, 0x5F, 0xFF},
	{0xAF, 0x5F, 0x87, 0xFF},
	{0xAF, 0x5F, 0xAF, 0xFF},
	{0xAF, 0x5F, 0xD7, 0xFF},
	{0xAF, 0x5F, 0xFF, 0xFF},
	{0xAF, 0x87, 0x00, 0xFF},
	{0xAF, 0x87, 0x5F, 0xFF},
	{0xAF, 0x87, 0x87, 0xFF},
	{0xAF, 0x87, 0xAF, 0xFF},
	{0xAF, 0x87, 0xD7, 0xFF},
	{0xAF, 0x87, 0xFF, 0xFF},
	{0xAF, 0xAF, 0x00, 0xFF},
	{0xAF, 0xAF, 0x5F, 0xFF},
	{0xAF, 0xAF, 0x87, 0xFF},
	{0xAF, 0xAF, 0xAF, 0xFF},
	{0xAF, 0xAF, 0xD7, 0xFF},
	{0xAF, 0xAF, 0xFF, 0xFF},
	{0xAF, 0xD7, 0x00, 0xFF},
	{0xAF, 0xD7, 0x5F, 0xFF},
	{0xAF, 0xD7, 0x87, 0xFF},
	{0xAF, 0xD7, 0xAF, 0xFF},
	{0xAF, 0xD7, 0xD7, 0xFF},
	{0xAF, 0xD7, 0xFF, 0xFF},
	{0xAF, 0xFF, 0x00, 0xFF},
	{0xAF, 0xFF, 0x5F, 0xFF},
	{0xAF, 0xFF, 0x87, 0xFF},
	{0xAF, 0xFF, 0xAF, 0xFF},
	{0xAF, 0xFF, 0xD7, 0xFF},
	{0xAF, 0xFF, 0xFF, 0xFF},
	{0xD7, 0x00, 0x00, 0xFF},
	{0xD7, 0x00, 0x5F, 0xFF},
	{0xD7, 0x00, 0x87, 0xFF},
	{0xD7, 0x00, 0xAF, 0xFF},
	{0xD7, 0x00, 0xD7, 0xFF},
	{0xD7, 0x00, 0xFF, 0xFF},
	{0xD7, 0x5F, 0x00, 0xFF},
	{0xD7, 0x5F, 0x5F, 0xFF},
	{0xD7, 0x5F, 0x87, 0xFF},
	{0xD7, 0x5F, 0xAF, 0xFF},
	{0xD7, 0x5F, 0xD7, 0xFF},
	{0xD7, 0x5F, 0xFF, 0xFF},
	{0xD7, 0x87, 0x00, 0xFF},
	{0xD7, 0x87, 0x5F, 0xFF},
	{0xD7, 0x87, 0x87, 0xFF},
	{0xD7, 0x87, 0xAF, 0xFF},
	{0xD7, 0x87, 0xD7, 0xFF},
	{0xD7, 0x87, 0xFF, 0xFF},
	{0xD7, 0xAF, 0x00, 0xFF},
	{0xD7, 0xAF, 0x5F, 0xFF},
	{0xD7, 0xAF, 0x87, 0xFF},
	{0xD7, 0xAF, 0xAF, 0xFF},
	{0xD7, 0xAF, 0xD7, 0xFF},
	{0xD7, 0xAF, 0xFF, 0xFF},
	{0xD7, 0xD7, 0x00, 0xFF},
	{0xD7, 0xD7, 0x5F, 0xFF},
	{0xD7, 0xD7, 0x87, 0xFF},
	{0xD7, 0xD7, 0xAF, 0xFF},
	{0xD7, 0xD7, 0xD7, 0xFF},
	{0xD7, 0xD7, 0xFF, 0xFF},
	{0xD7, 0xFF, 0x00, 0xFF},
	{0xD7, 0xFF, 0x5F, 0xFF},
	{0xD7, 0xFF, 0x87, 0xFF},
	{0xD7, 0xFF, 0xAF, 0xFF},
	{0xD7, 0xFF, 0xD7, 0xFF},
	{0xD7, 0xFF, 0xFF, 0xFF},
	{0xFF, 0x00, 0x00, 0xFF},
	{0xFF, 0x00, 0x5F, 0xFF},
	{0xFF, 0x00, 0x87, 0xFF},
	{0xFF, 0x00, 0xAF, 0xFF},
	{0xFF, 0x00, 0xD7, 0xFF},
	{0xFF, 0x00, 0xFF, 0xFF},
	{0xFF, 0x5F, 0x00, 0xFF},
	{0xFF, 0x5F, 0x5F, 0xFF},
	{0xFF, 0x5F, 0x87, 0xFF},
	{0xFF, 0x5F, 0xAF, 0xFF},
	{0xFF, 0x5F, 0xD7, 0xFF},
	{0xFF, 0x5F, 0xFF, 0xFF},
	{0xFF, 0x87, 0x00, 0xFF},
	{0xFF, 0x87, 0x5F, 0xFF},
	{0xFF, 0x87, 0x87, 0xFF},
	{0xFF, 0x87, 0xAF, 0xFF},
	{0xFF, 0x87, 0xD7, 0xFF},
	{0xFF, 0x87, 0xFF, 0xFF},
	{0xFF, 0xAF, 0x00, 0xFF},
	{0xFF, 0xAF, 0x5F, 0xFF},
	{0xFF, 0xAF, 0x87, 0xFF},
	{0xFF, 0xAF, 0xAF, 0xFF},
	{0xFF, 0xAF, 0xD7, 0xFF},
	{0xFF, 0xAF, 0xFF, 0xFF},
	{0xFF, 0xD7, 0x00, 0xFF},
	{0xFF, 0xD7, 0x5F, 0xFF},
	{0xFF, 0xD7, 0x87, 0xFF},
	{0xFF, 0xD7, 0xAF, 0xFF},
	{0xFF, 0xD7, 0xD7, 0xFF},
	{0xFF, 0xD7, 0xFF, 0xFF},
	{0xFF, 0xFF, 0x00, 0xFF},
	{0xFF, 0xFF, 0x5F, 0xFF},
	{0xFF, 0xFF, 0x87, 0xFF},
	{0xFF, 0xFF, 0xAF, 0xFF},
	{0xFF, 0xFF, 0xD7, 0xFF},
	{0xFF, 0xFF, 0xFF, 0xFF},
	{0x08, 0x08, 0x08, 0xFF},
	{0x12, 0x12, 0x12, 0xFF},
	{0x1C, 0x1C, 0x1C, 0xFF},
	{0x26, 0x26, 0x26, 0xFF},
	{0x30, 0x30, 0x30, 0xFF},
	{0x3A, 0x3A, 0x3A, 0xFF},
	{0x44, 0x44, 0x44, 0xFF},
	{0x4E, 0x4E, 0x4E, 0xFF},
	{0x58, 0x58, 0x58, 0xFF},
	{0x60, 0x60, 0x60, 0xFF},
	{0x66, 0x66, 0x66, 0xFF},
	{0x76, 0x76, 0x76, 0xFF},
	{0x80, 0x80, 0x80, 0xFF},
	{0x8A, 0x8A, 0x8A, 0xFF},
	{0x94, 0x94, 0x94, 0xFF},
	{0x9E, 0x9E, 0x9E, 0xFF},
	{0xA8, 0xA8, 0xA8, 0xFF},
	{0xB2, 0xB2, 0xB2, 0xFF},
	{0xBC, 0xBC, 0xBC, 0xFF},
	{0xC6, 0xC6, 0xC6, 0xFF},
	{0xD0, 0xD0, 0xD0, 0xFF},
	{0xDA, 0xDA, 0xDA, 0xFF},
	{0xE4, 0xE4, 0xE4, 0xFF},
	{0xEE, 0xEE, 0xEE, 0xFF},
}

// Pick the closest term-256 color.
//
// The color is converted to alpha-multiplied RGBA so this is important that
// alpha is set to 255 for full color. The return value is between 0 and 255.
func rgbToANSI(c color.NRGBA) int {
	// "Optimized" version of color.Palette{}.Index() that premultiplies c.A then
	// discards c.A.
	closest := 0
	delta := 1<<31 - 1
	cR, cG, cB, _ := c.RGBA()
	r := int(cR >> 8)
	g := int(cG >> 8)
	b := int(cB >> 8)
	for i, col := range lookupANSIToRGB {
		dR := r - int(col.R)
		dG := g - int(col.G)
		dB := b - int(col.B)
		d := dR*dR + dG*dG + dB*dB
		if d < delta {
			delta = d
			closest = i
		}
	}
	return closest
}
