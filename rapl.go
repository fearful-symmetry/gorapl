package gorapl

import (
	"fmt"

	"github.com/fearful-symmetry/gomsr"
)

//RAPLHandler manages a stateful connection to the RAPL system.
type RAPLHandler struct {
	availDomains []RAPLDomain //Available RAPL domains
	domainMask   uint         //a bitmask to make it easier to find available domains
	msrDev       gomsr.MSRDev
	units        RAPLPowerUnit
}

//NewRAPL returns a new RAPL handler
func NewRAPL(cpu int) (RAPLHandler, error) {

	//TODO: eventually we'll need to handle multiple CPU packages

	domains, mask := getAvailableDomains(cpu)
	if len(domains) == 0 {
		return RAPLHandler{}, fmt.Errorf("No RAPL domains available on CPU")
	}

	msr, err := gomsr.MSR(cpu)
	if err != nil {
		return RAPLHandler{}, err
	}

	handler := RAPLHandler{availDomains: domains, domainMask: mask, msrDev: msr}

	handler.units, err = handler.ReadPowerUnit()
	if err != nil {
		return RAPLHandler{}, err
	}

	return handler, nil
}

//ReadPowerLimit returns the MSR_[DOMAIN]_POWER_LIMIT MSR
func (h RAPLHandler) ReadPowerLimit(domain RAPLDomain) (RAPLPowerLimit, error) {

	if (domain.mask & h.domainMask) == 0 {
		return RAPLPowerLimit{}, fmt.Errorf("Domain %s does not exist on system", domain.name)
	}

	data, err := h.msrDev.Read(domain.msrs.PowerLimit)
	if err != nil {
		return RAPLPowerLimit{}, err
	}
	return parsePowerLimit(data), nil
}

//ReadPowerUnit returns the MSR_RAPL_POWER_UNIT MSR
func (h RAPLHandler) ReadPowerUnit() (RAPLPowerUnit, error) {

	data, err := h.msrDev.Read(MSRPowerUnit)
	if err != nil {
		return RAPLPowerUnit{}, err
	}

	return parsePowerUnit(data), nil

}

// helper functions

//Borrowed this from the kernel. Traverse over the Energy Status MSRs to see what RAPL domains are available
func getAvailableDomains(cpu int) ([]RAPLDomain, uint) {

	var availDomains []RAPLDomain
	var dm uint

	if _, exists := gomsr.ReadMSR(cpu, Package.msrs.EnergyStatus); exists == nil {
		availDomains = append(availDomains, Package)
		dm = dm | Package.mask
	}

	if _, exists := gomsr.ReadMSR(cpu, DRAM.msrs.EnergyStatus); exists == nil {
		availDomains = append(availDomains, DRAM)
		dm = dm | DRAM.mask
	}

	if _, exists := gomsr.ReadMSR(cpu, PP0.msrs.Policy); exists == nil {
		availDomains = append(availDomains, PP0)
		dm = dm | PP0.mask
	}

	if _, exists := gomsr.ReadMSR(cpu, PP1.msrs.EnergyStatus); exists == nil {
		availDomains = append(availDomains, PP1)
		dm = dm | PP1.mask
	}

	return availDomains, dm
}
