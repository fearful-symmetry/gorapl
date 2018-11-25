package gorapl

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPowerLimit(t *testing.T) {

	h, err := NewRAPL(0)
	if err != nil {
		t.Fatalf("Could not init: %s", err)
	}
	dat, err := h.ReadPowerLimit(DRAM)
	if err != nil {
		t.Fatalf("Could not read MSR: %s", err)
	}
	fmt.Printf("Got back: %#v\n", dat)

}

func TestParsePowerLimit(t *testing.T) {
	var msrVal uint64 = 0xf0000140c00
	knownParsed := RAPLPowerLimit{Limit1: PowerLimitSetting{
		PowerLimit:      0xc00,
		EnableLimit:     false,
		ClampingLimit:   false,
		TimeWindowLimit: 0xa,
	},
		Limit2: PowerLimitSetting{
			PowerLimit:      0xf00,
			EnableLimit:     false,
			ClampingLimit:   false,
			TimeWindowLimit: 0x0,
		},
		Lock: false,
	}

	parsedMsr := parsePowerLimit(msrVal)

	if !reflect.DeepEqual(parsedMsr, knownParsed) {
		t.Fatalf("struct failed: %#v", parsedMsr)
	}

	//fmt.Printf("%#v\n", parsedMsr)
}

func TestParsePowerUnit(t *testing.T) {
	var msrVal uint64 = 0xa0e03
	knownParsed := RAPLPowerUnit{
		PowerUnits:        0x3,
		EnergyStatusUnits: 0xe,
		TimeUnits:         0xa,
	}
	parsedMSr := parsePowerUnit(msrVal)

	if !reflect.DeepEqual(parsedMSr, knownParsed) {
		t.Fatalf("struct failed: %#v", parsedMSr)
	}
}
