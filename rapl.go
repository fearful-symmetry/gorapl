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
//This MSR defines power limits for the given domain. Every domain has this MSR
func (h RAPLHandler) ReadPowerLimit(domain RAPLDomain) (RAPLPowerLimit, error) {

	if (domain.mask & h.domainMask) == 0 {
		return RAPLPowerLimit{}, fmt.Errorf("Domain %s does not exist on system", domain.name)
	}

	data, err := h.msrDev.Read(domain.msrs.PowerLimit)
	if err != nil {
		return RAPLPowerLimit{}, err
	}

	var singleLimit = false
	if domain != Package {
		singleLimit = true
	}

	return parsePowerLimit(data, h.units, singleLimit), nil
}

//ReadEnergyStatus returns the MSR_[DOMAIN]_ENERGY_STATUS MSR
//This MSR is a single 32 bit field that reports the energy usage for the domain.
//Updated ~1ms. Every domain has this MSR. This is a cumulative register
func (h RAPLHandler) ReadEnergyStatus(domain RAPLDomain) (float64, error) {

	if (domain.mask & h.domainMask) == 0 {
		return 0, fmt.Errorf("Domain %s does not exist on system", domain.name)
	}

	data, err := h.msrDev.Read(domain.msrs.EnergyStatus)
	if err != nil {
		return 0, err
	}

	return float64(data&0xffffffff) * h.units.EnergyStatusUnits, nil

}

//ReadPolicy returns the MSR_[DOMAIN]_POLICY msr. This constists of a single value.
//The value is a priority that balances energy between the core and uncore devices. It's only available on the PP0/PP1 domains.
func (h RAPLHandler) ReadPolicy(domain RAPLDomain) (uint64, error) {

	if (domain.mask & h.domainMask) == 0 {
		return 0, fmt.Errorf("Domain %s does not exist on system", domain.name)
	}

	if domain.msrs.Policy == 0 {
		return 0, fmt.Errorf("Domain %s does not support the POLICY MSR", domain.name)
	}

	data, err := h.msrDev.Read(domain.msrs.Policy)
	if err != nil {
		return 0, err
	}

	return data & 0x1f, nil

}

//ReadPerfStatus returns the MSR_[DOMAIN]_PERF_STATUS msr. This is a single value.
//The value is the amount of time that the domain has been throttled due to RAPL limits. This is not available on PP1.
func (h RAPLHandler) ReadPerfStatus(domain RAPLDomain) (float64, error) {

	if (domain.mask & h.domainMask) == 0 {
		return 0, fmt.Errorf("Domain %s does not exist on system", domain.name)
	}

	if domain.msrs.PerfStatus == 0 {
		return 0, fmt.Errorf("Domain %s does not support the POLICY MSR", domain.name)
	}

	data, err := h.msrDev.Read(domain.msrs.PerfStatus)
	if err != nil {
		return 0, err
	}

	return float64(data&0xffffffff) * h.units.TimeUnits, nil
}

//ReadPowerInfo returns the MSR_[DOMAIN]_POWER_INFO MSR. This MSR is not available on PP0/PP1
func (h RAPLHandler) ReadPowerInfo(domain RAPLDomain) (RAPLPowerInfo, error) {

	if (domain.mask & h.domainMask) == 0 {
		return RAPLPowerInfo{}, fmt.Errorf("Domain %s does not exist on system", domain.name)
	}

	if domain.msrs.PerfStatus == 0 {
		return RAPLPowerInfo{}, fmt.Errorf("Domain %s does not support the POLICY MSR", domain.name)
	}

	data, err := h.msrDev.Read(domain.msrs.PowerInfo)
	if err != nil {
		return RAPLPowerInfo{}, err
	}

	return parsePowerInfo(data, h.units), nil
}

//ReadPowerUnit returns the MSR_RAPL_POWER_UNIT MSR
//This has no associated domain
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
