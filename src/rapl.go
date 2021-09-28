package gorapl

import (
	"errors"
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

//RAPLOptions is used to set options for a new RAPLHandler
type RAPLOptions struct {
	FmtPattern string
	CPU        int
}

// ErrMSRDoesNotExist is the error for instances when a Domain does not exist on a given RAPL domain
var ErrMSRDoesNotExist = errors.New("MSR does not exist on selected Domain")

//CreateNewHandler creates a RAPL register handler for the given CPU
func CreateNewHandler(cpu int, fmtS string) (RAPLHandler, error) {

	domains, mask := getAvailableDomains(cpu)
	if len(domains) == 0 {
		return RAPLHandler{}, fmt.Errorf("No RAPL domains available on CPU")
	}

	var msr gomsr.MSRDev
	var err error
	if fmtS == "" {
		msr, err = gomsr.MSR(cpu)
		if err != nil {
			return RAPLHandler{}, err
		}
	} else {
		msr, err = gomsr.MSRWithLocation(cpu, fmtS)
		if err != nil {
			return RAPLHandler{}, err
		}
	}

	handler := RAPLHandler{availDomains: domains, domainMask: mask, msrDev: msr}

	handler.units, err = handler.ReadPowerUnit()
	if err != nil {
		return RAPLHandler{}, err
	}

	return handler, nil
}

//NewRAPLWithOptions returns a new handler with the given options.
//The CPU setting determines the logical processor opened via the given fmt pattern.
//The default pattern is /dev/cpu/%d/msr_safe.
func NewRAPLWithOptions(options RAPLOptions) (RAPLHandler, error) {

	return CreateNewHandler(options.CPU, options.FmtPattern)

}

//NewRAPL returns a new RAPL handler
func NewRAPL() (RAPLHandler, error) {

	//TODO: eventually we'll need to handle multiple CPU packages

	return CreateNewHandler(0, "")
}

// GetDomains returns the list of RAPL domains on the package
func (h RAPLHandler) GetDomains() []RAPLDomain {
	return h.availDomains
}

//ReadPowerLimit returns the MSR_[DOMAIN]_POWER_LIMIT MSR
//This MSR defines power limits for the given domain. Every domain has this MSR
func (h RAPLHandler) ReadPowerLimit(domain RAPLDomain) (RAPLPowerLimit, error) {

	if (domain.Mask & h.domainMask) == 0 {
		return RAPLPowerLimit{}, fmt.Errorf("Domain %s does not exist on system", domain.Name)
	}

	data, err := h.msrDev.Read(domain.MSRs.PowerLimit)
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

	if (domain.Mask & h.domainMask) == 0 {
		return 0, fmt.Errorf("Domain %s does not exist on system", domain.Name)
	}

	data, err := h.msrDev.Read(domain.MSRs.EnergyStatus)
	if err != nil {
		return 0, err
	}

	return float64(data&0xffffffff) * h.units.EnergyStatusUnits, nil

}

//ReadPolicy returns the MSR_[DOMAIN]_POLICY msr. This constists of a single value.
//The value is a priority that balances energy between the core and uncore devices. It's only available on the PP0/PP1 domains.
func (h RAPLHandler) ReadPolicy(domain RAPLDomain) (uint64, error) {

	if (domain.Mask & h.domainMask) == 0 {
		return 0, fmt.Errorf("Domain %s does not exist on system", domain.Name)
	}

	if domain.MSRs.Policy == 0 {
		return 0, ErrMSRDoesNotExist
	}

	data, err := h.msrDev.Read(domain.MSRs.Policy)
	if err != nil {
		return 0, err
	}

	return data & 0x1f, nil

}

//ReadPerfStatus returns the MSR_[DOMAIN]_PERF_STATUS msr. This is a single value.
//The value is the amount of time that the domain has been throttled due to RAPL limits. This is not available on PP1.
func (h RAPLHandler) ReadPerfStatus(domain RAPLDomain) (float64, error) {

	if (domain.Mask & h.domainMask) == 0 {
		return 0, fmt.Errorf("Domain %s does not exist on system", domain.Name)
	}

	if domain.MSRs.PerfStatus == 0 {
		return 0, ErrMSRDoesNotExist
	}

	data, err := h.msrDev.Read(domain.MSRs.PerfStatus)
	if err != nil {
		return 0, err
	}

	return float64(data&0xffffffff) * h.units.TimeUnits, nil
}

//ReadPowerInfo returns the MSR_[DOMAIN]_POWER_INFO MSR. This MSR is not available on PP0/PP1
func (h RAPLHandler) ReadPowerInfo(domain RAPLDomain) (RAPLPowerInfo, error) {

	if (domain.Mask & h.domainMask) == 0 {
		return RAPLPowerInfo{}, fmt.Errorf("Domain %s does not exist on system", domain.Name)
	}

	if domain.MSRs.PerfStatus == 0 {
		return RAPLPowerInfo{}, ErrMSRDoesNotExist
	}

	data, err := h.msrDev.Read(domain.MSRs.PowerInfo)
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

	if _, exists := gomsr.ReadMSR(cpu, Package.MSRs.EnergyStatus); exists == nil {
		availDomains = append(availDomains, Package)
		dm = dm | Package.Mask
	}

	if _, exists := gomsr.ReadMSR(cpu, DRAM.MSRs.EnergyStatus); exists == nil {
		availDomains = append(availDomains, DRAM)
		dm = dm | DRAM.Mask
	}

	if _, exists := gomsr.ReadMSR(cpu, PP0.MSRs.Policy); exists == nil {
		availDomains = append(availDomains, PP0)
		dm = dm | PP0.Mask
	}

	if _, exists := gomsr.ReadMSR(cpu, PP1.MSRs.EnergyStatus); exists == nil {
		availDomains = append(availDomains, PP1)
		dm = dm | PP1.Mask
	}

	return availDomains, dm
}
