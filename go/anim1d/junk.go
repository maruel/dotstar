// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// These are incomplete and will either be removed or fixed.

package anim1d

import (
	"math/rand"
	"time"
)

// TODO(maruel): Create NightSky with:
//    - Stars
//    - WishingStar
//    - Aurores
//    - Super nova.
//    - Rotation de la terre?
//    - Station Internationale?

// Aurore commence lentement, se transforme lentement et éventuellement
// disparait.
type Aurore struct {
}

func (a *Aurore) NextFrame(pixels Frame, timeMS uint32) {
	// TODO(maruel): Redo.
	y := float32(timeMS) * .01
	for i := range pixels {
		x := float32(i)
		//a := 32 + 31*sin(x/(37.+15*cos(y/74)))*cos(y/(31+11*sin(x/57)))
		b := (32 + 31*(sin(hypot(200-y, 320-x)/16))) * (0.5 + 0.5*sin(y*0.1))
		pixels[i].R = 0
		//pixels[i].G = uint8(a + b)
		pixels[i].G = uint8(b)
		pixels[i].B = 0
	}
}

type star struct {
	intensity uint8
	// TODO(maruel): Use maruel/temperature.
}

type NightStars struct {
	stars []star
}

func (e *NightStars) NextFrame(pixels Frame, timeMS uint32) {
	if len(e.stars) != len(pixels) {
		r := rand.NewSource(0)
		e.stars = make([]star, len(pixels))
		for i := range e.stars {
			j := r.Int63()
			// Cut off at 25%.
			if j&0x30000 != 0x30000 {
				continue
			}
			// Use gamma == 2 and limit intensity at 50%.
			d := int(j&0xff+1) * int((j>>8)&0xff+1)
			e.stars[i].intensity = uint8((d-1)>>8) / 2
		}
	}

	r := rand.NewSource(int64((&Rand{}).Eval(timeMS)))
	for i, s := range e.stars {
		if s.intensity == 0 {
			pixels[i] = Color{}
			continue
		}
		j := r.Int63()
		// Use gamma == 2.
		d := int(j&0xf+1) * int((j>>4)&0xf+1)
		y := (d-1)>>4 + int(s.intensity)
		if y > 255 {
			y = 255
		} else if y < 0 {
			y = 0
		}
		f := uint8(y)
		pixels[i] = Color{f, f, f}
	}
}

// WishingStar draws a wishing star from time to time.
//
// It will only draw one star at a time. To increase the likelihood of getting
// many simultaneously, create multiple instances.
type WishingStar struct {
	Duration     time.Duration // Average duration of a star.
	AverageDelay time.Duration // Average delay between each wishing star.
}

func (w *WishingStar) NextFrame(pixels Frame, timeMS uint32) {
	/*
		// Create a deterministic replay by using the current number of
		// the wishing star as the seed for the current flow. Make it independent of
		// any other non-deterministic source.
		i := timeMS / w.AverageDelay
		r := rand.New(rand.NewSource(int64(i)))
		// Always calculate things in the same order to keep the calculation
		// deterministic.
		startOffset := r.Int63()
		startPos := r.Int63()
		intensity := r.Int63()
		orientation := r.Intn(2)
		// Draw according to these parameters.
		// - Trail
		// - Observed speed based on orientation
	*/
}
