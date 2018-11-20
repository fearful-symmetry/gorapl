package gorapl

import (
	"fmt"
	"github.com/fearful-symmetry/gomsr"
)

//RAPLHandler manages a stateful connection to the RAPL system.
type RAPLHandler struct {
	availDomains []RAPLDomain //Available RAPL domains
	msrDev       gomsr.MSRDev
}

//NewRAPL returns a new RAPL handler
func NewRAPL() (RAPLHandler, error) {

	//TODO: eventually we'll need to handle multiple CPU packages
	cpu := 0

	domains := getAvailableDomains(cpu)
	if len(domains) == 0 {
		return RAPLHandler{}, fmt.Errorf("No RAPL domains available on CPU")
	}

	msr, err := gomsr.MSR(cpu)
	if err != nil {
		return RAPLHandler{}, err
	}

	return RAPLHandler{availDomains: domains, msrDev: msr}, nil
}

//Borrowed this from the kernel. Traverse over the Energy Status MSRs to see what RAPL domains are available
func getAvailableDomains(cpu int) []RAPLDomain {

	var availDomains []RAPLDomain

	if _, exists := gomsr.ReadMSR(cpu, globalMSR.Pkg.EnergyStatus); exists == nil {
		availDomains = append(availDomains, Package)
	}

	if _, exists := gomsr.ReadMSR(cpu, globalMSR.DRAM.EnergyStatus); exists == nil {
		availDomains = append(availDomains, DRAM)
	}

	if _, exists := gomsr.ReadMSR(cpu, globalMSR.PP0.Policy); exists == nil {
		availDomains = append(availDomains, PP0)
	}

	if _, exists := gomsr.ReadMSR(cpu, globalMSR.PP1.EnergyStatus); exists == nil {
		availDomains = append(availDomains, PP1)
	}

	return availDomains
}
