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
	var msrVal uint64 = 0x7fd00014ea82
	knownParsed := RAPLPowerLimit{Limit1: PowerLimitSetting{
		PowerLimit:      3408.25,
		EnableLimit:     true,
		ClampingLimit:   false,
		TimeWindowLimit: 1,
	},
		Limit2: PowerLimitSetting{
			PowerLimit:      4090,
			EnableLimit:     false,
			ClampingLimit:   false,
			TimeWindowLimit: 0.0009765625,
		},
		Lock: false,
	}

	units := parsePowerUnit(0xa0e03)

	parsedMsr := parsePowerLimit(msrVal, units, false)
	fmt.Printf("%#v\n", parsedMsr)
	if !reflect.DeepEqual(parsedMsr, knownParsed) {
		t.Fatalf("struct failed: %#v", parsedMsr)
	}

	//fmt.Printf("%#v\n", parsedMsr)
}

func TestParsePowerUnit(t *testing.T) {
	var msrVal uint64 = 0xa0e03
	knownParsed := RAPLPowerUnit{
		PowerUnits:        0.125,
		EnergyStatusUnits: 6.103515625e-05,
		TimeUnits:         0.0009765625,
	}
	parsedMSr := parsePowerUnit(msrVal)
	//fmt.Printf("%#v\n", parsedMSr)
	if !reflect.DeepEqual(parsedMSr, knownParsed) {
		t.Fatalf("struct failed: %#v", parsedMSr)
	}
}

func TestParsePowerInfo(t *testing.T) {
	var msrVal uint64 = 0x2a0
	knownParsed := RAPLPowerInfo{
		ThermalSpecPower: 84,
		MinPower:         0x0,
		MaxPower:         0x0,
	}
	units := RAPLPowerUnit{
		PowerUnits:        0.125,
		EnergyStatusUnits: 6.103515625e-05,
		TimeUnits:         0.0009765625,
	}
	parsedMSr := parsePowerInfo(msrVal, units)
	//fmt.Printf("%#v\n", parsedMSr)
	if !reflect.DeepEqual(parsedMSr, knownParsed) {
		t.Fatalf("struct failed: %#v", parsedMSr)
	}
}
