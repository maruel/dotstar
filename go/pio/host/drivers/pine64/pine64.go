// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package pine64

import (
	"github.com/maruel/dlibox/go/pio/host/drivers/allwinner"
	"github.com/maruel/dlibox/go/pio/host/pins"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

// Version is the board version. Only reports as 1 for now.
var Version int = 1

var (
	P1_1  gpio.PinIO = pins.V3_3      // 3.3 volt; max 40mA
	P1_2  gpio.PinIO = pins.V5        // 5 volt (before filtering)
	P1_3  gpio.PinIO = allwinner.PH3  //
	P1_4  gpio.PinIO = pins.V5        //
	P1_5  gpio.PinIO = allwinner.PH2  //
	P1_6  gpio.PinIO = pins.GROUND    //
	P1_7  gpio.PinIO = allwinner.PL10 //
	P1_8  gpio.PinIO = allwinner.PB0  //
	P1_9  gpio.PinIO = pins.GROUND    //
	P1_10 gpio.PinIO = allwinner.PB1  //
	P1_11 gpio.PinIO = allwinner.PC7  //
	P1_12 gpio.PinIO = allwinner.PC8  //
	P1_13 gpio.PinIO = allwinner.PH9  //
	P1_14 gpio.PinIO = pins.GROUND    //
	P1_15 gpio.PinIO = allwinner.PC12 //
	P1_16 gpio.PinIO = allwinner.PC13 //
	P1_17 gpio.PinIO = pins.V3_3      //
	P1_18 gpio.PinIO = allwinner.PC14 //
	P1_19 gpio.PinIO = allwinner.PC0  //
	P1_20 gpio.PinIO = pins.GROUND    //
	P1_21 gpio.PinIO = allwinner.PC1  //
	P1_22 gpio.PinIO = allwinner.PC15 //
	P1_23 gpio.PinIO = allwinner.PC2  //
	P1_24 gpio.PinIO = allwinner.PC3  //
	P1_25 gpio.PinIO = pins.GROUND    //
	P1_26 gpio.PinIO = allwinner.PH7  //
	P1_27 gpio.PinIO = allwinner.PL9  //
	P1_28 gpio.PinIO = allwinner.PL8  //
	P1_29 gpio.PinIO = allwinner.PH5  //
	P1_30 gpio.PinIO = pins.GROUND    //
	P1_31 gpio.PinIO = allwinner.PH6  //
	P1_32 gpio.PinIO = allwinner.PC4  //
	P1_33 gpio.PinIO = allwinner.PC5  //
	P1_34 gpio.PinIO = pins.GROUND    //
	P1_35 gpio.PinIO = allwinner.PC9  //
	P1_36 gpio.PinIO = allwinner.PC6  //
	P1_37 gpio.PinIO = allwinner.PC16 //
	P1_38 gpio.PinIO = allwinner.PC10 //
	P1_39 gpio.PinIO = pins.GROUND    //
	P1_40 gpio.PinIO = allwinner.PC11 //

	EULER_1  gpio.PinIO = pins.V3_3        //
	EULER_2  gpio.PinIO = pins.DC_IN       //
	EULER_3  gpio.PinIO = pins.BAT_PLUS    //
	EULER_4  gpio.PinIO = pins.DC_IN       //
	EULER_5  gpio.PinIO = pins.TEMP_SENSOR //
	EULER_6  gpio.PinIO = pins.GROUND      //
	EULER_7  gpio.PinIO = pins.IR_RX       //
	EULER_8  gpio.PinIO = pins.V5          //
	EULER_9  gpio.PinIO = pins.GROUND      //
	EULER_10 gpio.PinIO = allwinner.PH8    //
	EULER_11 gpio.PinIO = allwinner.PB3    //
	EULER_12 gpio.PinIO = allwinner.PB4    //
	EULER_13 gpio.PinIO = allwinner.PB5    //
	EULER_14 gpio.PinIO = pins.GROUND      //
	EULER_15 gpio.PinIO = allwinner.PB6    //
	EULER_16 gpio.PinIO = allwinner.PB7    //
	EULER_17 gpio.PinIO = pins.V3_3        //
	EULER_18 gpio.PinIO = allwinner.PD4    //
	EULER_19 gpio.PinIO = allwinner.PD2    //
	EULER_20 gpio.PinIO = pins.GROUND      //
	EULER_21 gpio.PinIO = allwinner.PD3    //
	EULER_22 gpio.PinIO = allwinner.PD5    //
	EULER_23 gpio.PinIO = allwinner.PD1    //
	EULER_24 gpio.PinIO = allwinner.PD0    //
	EULER_25 gpio.PinIO = pins.GROUND      //
	EULER_26 gpio.PinIO = allwinner.PD6    //
	EULER_27 gpio.PinIO = allwinner.PB2    //
	EULER_28 gpio.PinIO = allwinner.PD7    //
	EULER_29 gpio.PinIO = allwinner.PB8    //
	EULER_30 gpio.PinIO = allwinner.PB9    //
	EULER_31 gpio.PinIO = pins.EAROUTP     //
	EULER_32 gpio.PinIO = pins.EAROUT_N    //
	EULER_33 gpio.PinIO = gpio.INVALID     //
	EULER_34 gpio.PinIO = pins.GROUND      //

	EXP_1  gpio.PinIO = pins.V3_3        //
	EXP_2  gpio.PinIO = allwinner.PL7    //
	EXP_3  gpio.PinIO = pins.CHARGER_LED //
	EXP_4  gpio.PinIO = pins.RESET       //
	EXP_5  gpio.PinIO = pins.PWR_SWITCH  //
	EXP_6  gpio.PinIO = pins.GROUND      //
	EXP_7  gpio.PinIO = allwinner.PB8    //
	EXP_8  gpio.PinIO = allwinner.PB9    //
	EXP_9  gpio.PinIO = pins.GROUND      //
	EXP_10 gpio.PinIO = pins.KEY_ADC     //

	WIFI_BT_1  gpio.PinIO = pins.GROUND    //
	WIFI_BT_2  gpio.PinIO = allwinner.PG6  //
	WIFI_BT_3  gpio.PinIO = allwinner.PG0  //
	WIFI_BT_4  gpio.PinIO = allwinner.PG7  //
	WIFI_BT_5  gpio.PinIO = pins.GROUND    //
	WIFI_BT_6  gpio.PinIO = allwinner.PG8  //
	WIFI_BT_7  gpio.PinIO = allwinner.PG1  //
	WIFI_BT_8  gpio.PinIO = allwinner.PG9  //
	WIFI_BT_9  gpio.PinIO = allwinner.PG2  //
	WIFI_BT_10 gpio.PinIO = allwinner.PG10 //
	WIFI_BT_11 gpio.PinIO = allwinner.PG3  //
	WIFI_BT_12 gpio.PinIO = allwinner.PG11 //
	WIFI_BT_13 gpio.PinIO = allwinner.PG4  //
	WIFI_BT_14 gpio.PinIO = allwinner.PG12 //
	WIFI_BT_15 gpio.PinIO = allwinner.PG5  //
	WIFI_BT_16 gpio.PinIO = allwinner.PG13 //
	WIFI_BT_17 gpio.PinIO = allwinner.PL2  //
	WIFI_BT_18 gpio.PinIO = pins.GROUND    //
	WIFI_BT_19 gpio.PinIO = allwinner.PL3  //
	WIFI_BT_20 gpio.PinIO = allwinner.PL5  //
	WIFI_BT_21 gpio.PinIO = pins.X32KFOUT  //
	WIFI_BT_22 gpio.PinIO = allwinner.PL5  //
	WIFI_BT_23 gpio.PinIO = pins.GROUND    //
	WIFI_BT_24 gpio.PinIO = allwinner.PL6  //
	WIFI_BT_25 gpio.PinIO = pins.VCC       //
	WIFI_BT_26 gpio.PinIO = pins.IOVCC     //

	AUDIO_LEFT  gpio.PinIO = gpio.INVALID // TODO(maruel): Figure out, is that EAROUT?
	AUDIO_RIGHT gpio.PinIO = gpio.INVALID //
)

// See headers.Headers for more info.
var Headers = map[string][][]gpio.PinIO{
	"P1": {
		{P1_1, P1_2},
		{P1_3, P1_4},
		{P1_5, P1_6},
		{P1_7, P1_8},
		{P1_9, P1_10},
		{P1_11, P1_12},
		{P1_13, P1_14},
		{P1_15, P1_16},
		{P1_17, P1_18},
		{P1_19, P1_20},
		{P1_21, P1_22},
		{P1_23, P1_24},
		{P1_25, P1_26},
		{P1_27, P1_28},
		{P1_29, P1_30},
		{P1_31, P1_32},
		{P1_33, P1_34},
		{P1_35, P1_36},
		{P1_37, P1_38},
		{P1_39, P1_20},
	},
	"EULER": {
		{EULER_1, EULER_2},
		{EULER_3, EULER_4},
		{EULER_5, EULER_6},
		{EULER_7, EULER_8},
		{EULER_9, EULER_10},
		{EULER_11, EULER_12},
		{EULER_13, EULER_14},
		{EULER_15, EULER_16},
		{EULER_17, EULER_18},
		{EULER_19, EULER_20},
		{EULER_21, EULER_22},
		{EULER_23, EULER_24},
		{EULER_25, EULER_26},
		{EULER_27, EULER_28},
		{EULER_29, EULER_30},
		{EULER_31, EULER_32},
		{EULER_33, EULER_34},
	},
	"EXP": {
		{EXP_1, EXP_2},
		{EXP_3, EXP_4},
		{EXP_5, EXP_6},
		{EXP_7, EXP_8},
		{EXP_9, EXP_10},
	},
	"WIFI_BT": {
		{WIFI_BT_1, WIFI_BT_2},
		{WIFI_BT_3, WIFI_BT_4},
		{WIFI_BT_5, WIFI_BT_6},
		{WIFI_BT_7, WIFI_BT_8},
		{WIFI_BT_9, WIFI_BT_10},
		{WIFI_BT_11, WIFI_BT_12},
		{WIFI_BT_13, WIFI_BT_14},
		{WIFI_BT_15, WIFI_BT_16},
		{WIFI_BT_17, WIFI_BT_18},
		{WIFI_BT_19, WIFI_BT_20},
		{WIFI_BT_21, WIFI_BT_22},
		{WIFI_BT_23, WIFI_BT_24},
		{WIFI_BT_25, WIFI_BT_26},
	},
	"AUDIO": {
		{AUDIO_LEFT},
		{AUDIO_RIGHT},
	},
}
