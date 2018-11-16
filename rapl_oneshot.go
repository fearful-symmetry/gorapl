package gorapl

import (
	"github.com/fearful-symmetry/gomsr"
)

//This file contains "one-shot" functions to quickly read from a given MSR without having to manage any kind of stateful objects.

/*
TODO:
Still unsure about the logic for accessing various MSR devices.
In theory, and based on what I've observed from intel's code,
the RAPL MSRs are package-wide, i.e it doesn't matter what MSR you access on
the package in question. Still, I'm not 100% sure.
*/

//ReadPowerLimit Reads the MSR_[DOMAIN]_POWER_LIMIT MSR for the given domain
func ReadPowerLimit(domain RAPLDomain) (RAPLPowerLimit, error) {

	msr, err := getDomainMSRs(domain)
	if err != nil {
		return RAPLPowerLimit{}, err
	}

	data, err := gomsr.ReadMSR(0, msr.PowerLimit)
	if err != nil {
		return RAPLPowerLimit{}, err
	}

	return parsePowerLimit(data), nil

}
