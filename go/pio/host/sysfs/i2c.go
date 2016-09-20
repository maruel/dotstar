// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package sysfs

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

// I2C is an open I²C bus via sysfs.
//
// It can be used to communicate with multiple devices from multiple goroutines.
type I2C struct {
	f  *os.File
	l  sync.Mutex // In theory the kernel probably has an internal lock but not taking any chance.
	fn functionality
}

func makeI2C(busNumber int) (*I2C, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", busNumber), os.O_RDWR, os.ModeExclusive)
	if err != nil {
		// Try to be helpful here. There are generally two cases:
		// - /dev/i2c-X doesn't exist. In this case, /boot/config.txt has to be
		//   edited to enable I²C then the device must be rebooted.
		// - permission denied. In this case, the user has to be added to plugdev.
		if os.IsNotExist(err) {
			return nil, errors.New("I²C is not configured; please follow instructions at https://github.com/maruel/dlibox/tree/master/go/setup")
		}
		return nil, fmt.Errorf("are you member of group 'plugdev'? please follow instructions at https://github.com/maruel/dlibox/tree/master/go/setup. %s", err)
	}
	i := &I2C{f: f}

	// TODO(maruel): Changing the speed is currently doing this for all devices.
	// https://github.com/raspberrypi/linux/issues/215
	// Need to access /sys/module/i2c_bcm2708/parameters/baudrate

	// Query to know if 10 bits addresses are supported.
	if err = i.ioctl(ioctlFuncs, uintptr(unsafe.Pointer(&i.fn))); err != nil {
		return nil, err
	}
	return i, nil
}

// Close closes the handle to the I²C driver. It is not a requirement to close
// before process termination.
func (i *I2C) Close() error {
	i.l.Lock()
	defer i.l.Unlock()
	err := i.f.Close()
	i.f = nil
	return err
}

// Tx execute a transaction as a single operation unit.
func (i *I2C) Tx(addr uint16, w, r []byte) error {
	if addr >= 0x400 || (addr >= 0x80 && i.fn&func10BIT_ADDR == 0) {
		return nil
	}

	// Convert the messages to the internal format.
	var buf [2]i2cMsg
	msgs := buf[:1]
	buf[0].addr = addr
	buf[0].length = uint16(len(w))
	buf[0].buf = uintptr(unsafe.Pointer(&w[0]))
	if len(r) != 0 {
		msgs = buf[:]
		buf[1].addr = addr
		buf[1].flags = flagRD
		buf[1].length = uint16(len(r))
		buf[1].buf = uintptr(unsafe.Pointer(&r[0]))
	}
	p := rdwrIoctlData{
		msgs:  uintptr(unsafe.Pointer(&msgs[0])),
		nmsgs: uint32(len(msgs)),
	}
	pp := uintptr(unsafe.Pointer(&p))
	i.l.Lock()
	defer i.l.Unlock()
	return i.ioctl(ioctlRdwr, pp)
}

func (i *I2C) ioctl(op uint, arg uintptr) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, i.f.Fd(), uintptr(op), arg); errno != 0 {
		return fmt.Errorf("i²c ioctl: %s", syscall.Errno(errno))
	}
	return nil
}

// Private details.

// i2cdev driver IOCTL control codes.
//
// Constants and structure definition can be found at
// /usr/include/linux/i2c-dev.h and /usr/include/linux/i2c.h.
const (
	ioctlRetries = 0x701 // TODO(maruel): Expose this
	ioctlTimeout = 0x702 // TODO(maruel): Expose this; in units of 10ms
	ioctlSlave   = 0x703
	ioctlTenBits = 0x704 // TODO(maruel): Expose this but the header says it's broken (!?)
	ioctlFuncs   = 0x705
	ioctlRdwr    = 0x707
)

// flags
const (
	flagTEN          = 0x0010 // this is a ten bit chip address
	flagRD           = 0x0001 // read data, from slave to master
	flagSTOP         = 0x8000 // if I2C_FUNC_PROTOCOL_MANGLING
	flagNOSTART      = 0x4000 // if I2C_FUNC_NOSTART
	flagREV_DIR_ADDR = 0x2000 // if I2C_FUNC_PROTOCOL_MANGLING
	flagIGNORE_NAK   = 0x1000 // if I2C_FUNC_PROTOCOL_MANGLING
	flagNO_RD_ACK    = 0x0800 // if I2C_FUNC_PROTOCOL_MANGLING
	flagRECV_LEN     = 0x0400 // length will be first received byte

)

type functionality uint64

const (
	funcI2C                    = 0x00000001
	func10BIT_ADDR             = 0x00000002
	funcPROTOCOL_MANGLING      = 0x00000004 // I2C_M_IGNORE_NAK etc.
	funcSMBUS_PEC              = 0x00000008
	funcNOSTART                = 0x00000010 // I2C_M_NOSTART
	funcSMBUS_BLOCK_PROC_CALL  = 0x00008000 // SMBus 2.0
	funcSMBUS_QUICK            = 0x00010000
	funcSMBUS_READ_BYTE        = 0x00020000
	funcSMBUS_WRITE_BYTE       = 0x00040000
	funcSMBUS_READ_BYTE_DATA   = 0x00080000
	funcSMBUS_WRITE_BYTE_DATA  = 0x00100000
	funcSMBUS_READ_WORD_DATA   = 0x00200000
	funcSMBUS_WRITE_WORD_DATA  = 0x00400000
	funcSMBUS_PROC_CALL        = 0x00800000
	funcSMBUS_READ_BLOCK_DATA  = 0x01000000
	funcSMBUS_WRITE_BLOCK_DATA = 0x02000000
	funcSMBUS_READ_I2C_BLOCK   = 0x04000000 // I2C-like block xfer
	funcSMBUS_WRITE_I2C_BLOCK  = 0x08000000 // w/ 1-byte reg. addr.
)

func (f functionality) String() string {
	var out []string
	if f&funcI2C != 0 {
		out = append(out, "I2C")
	}
	if f&func10BIT_ADDR != 0 {
		out = append(out, "10BIT_ADDR")
	}
	if f&funcPROTOCOL_MANGLING != 0 {
		out = append(out, "PROTOCOL_MANGLING")
	}
	if f&funcSMBUS_PEC != 0 {
		out = append(out, "SMBUS_PEC")
	}
	if f&funcNOSTART != 0 {
		out = append(out, "NOSTART")
	}
	if f&funcSMBUS_BLOCK_PROC_CALL != 0 {
		out = append(out, "SMBUS_BLOCK_PROC_CALL")
	}
	if f&funcSMBUS_QUICK != 0 {
		out = append(out, "SMBUS_QUICK")
	}
	if f&funcSMBUS_READ_BYTE != 0 {
		out = append(out, "SMBUS_READ_BYTE")
	}
	if f&funcSMBUS_WRITE_BYTE != 0 {
		out = append(out, "SMBUS_WRITE_BYTE")
	}
	if f&funcSMBUS_READ_BYTE_DATA != 0 {
		out = append(out, "SMBUS_READ_BYTE_DATA")
	}
	if f&funcSMBUS_WRITE_BYTE_DATA != 0 {
		out = append(out, "SMBUS_WRITE_BYTE_DATA")
	}
	if f&funcSMBUS_READ_WORD_DATA != 0 {
		out = append(out, "SMBUS_READ_WORD_DATA")
	}
	if f&funcSMBUS_WRITE_WORD_DATA != 0 {
		out = append(out, "SMBUS_WRITE_WORD_DATA")
	}
	if f&funcSMBUS_PROC_CALL != 0 {
		out = append(out, "SMBUS_PROC_CALL")
	}
	if f&funcSMBUS_READ_BLOCK_DATA != 0 {
		out = append(out, "SMBUS_READ_BLOCK_DATA")
	}
	if f&funcSMBUS_WRITE_BLOCK_DATA != 0 {
		out = append(out, "SMBUS_WRITE_BLOCK_DATA")
	}
	if f&funcSMBUS_READ_I2C_BLOCK != 0 {
		out = append(out, "SMBUS_READ_I2C_BLOCK")
	}
	if f&funcSMBUS_WRITE_I2C_BLOCK != 0 {
		out = append(out, "SMBUS_WRITE_I2C_BLOCK")
	}
	return strings.Join(out, "|")
}

type rdwrIoctlData struct {
	msgs  uintptr // Pointer to i2cMsg
	nmsgs uint32
}

type i2cMsg struct {
	addr   uint16 // Address to communicate with
	flags  uint16 // 1 for read, see i2c.h for more details
	length uint16
	buf    uintptr
}
