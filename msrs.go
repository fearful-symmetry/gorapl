package gorapl

// //PkgPowerLimit is the MSR for package-level power limits
// //Fun fact: you can set this limit to 0, and the processor, will still work.
// const pkgPowerLimit = 0x610

// //PkgEnergyStatus is a read-only MSR for reading the actual emergy use from the package
// const PkgEnergyStatus = 0x611

// //PkgPerfStatus	is a read-only MSR that reports the time the package was throttled due to RAPL limits
// const PkgPerfStatus = 0x613

// //PkgPowerInfo is a read-only MSR that reports power range info for the package
// const PkgPowerInfo = 0x14

// //MSRDRAMPowerLimit is the MSR for DRAM-domain power limits
// const MSRDRAMPowerLimit = 0x618

type domainMSRs struct {
	PowerLimit   int64
	EnergyStatus int64
	Policy       int64
	PerfStatus   int64
	PowerInfo    int64
}

//This is a somewhat unorthidox syntax, but I'm not really a fan of cramming packages full of global `const` objects.
//For now we'll use 0x0 for MSRs that undefined for a given domain
// var globalMSR = struct {
// 	Pkg  domainMSRs
// 	DRAM domainMSRs
// 	PP0  domainMSRs
// 	PP1  domainMSRs
// }{
// 	domainMSRs{0x610, 0x611, 0x0, 0x613, 0x614},
// 	domainMSRs{0x618, 0x619, 0x0, 0x61b, 0x61c},
// 	domainMSRs{0x638, 0x639, 0x63a, 0x63a, 0x0},
// 	domainMSRs{0x640, 0x641, 0x642, 0x0, 0x0},
// }

// The various RAPL domains

//RAPLDomain is a string type that covers the various RAPL domains
type RAPLDomain struct {
	mask uint
	name string
	msrs domainMSRs
}

//Package is the RAPL domain for the CPU package
var Package = RAPLDomain{0x1, "Package", domainMSRs{0x610, 0x611, 0x0, 0x613, 0x614}}

//DRAM is the RAPL domain for the DRAM
var DRAM = RAPLDomain{0x2, "DRAM", domainMSRs{0x618, 0x619, 0x0, 0x61b, 0x61c}}

//PP0 is the RAPL domain for the processor core
var PP0 = RAPLDomain{0x4, "PP0", domainMSRs{0x638, 0x639, 0x63a, 0x63a, 0x0}}

//PP1 is platform-dependant, although it usually referrs to some uncore power plane
var PP1 = RAPLDomain{0x8, "PP1", domainMSRs{0x640, 0x641, 0x642, 0x0, 0x0}}

// func getDomainMSRs(domain RAPLDomain) (domainMSRs, error) {

// 	switch domain {
// 	case Package:
// 		return globalMSR.Pkg, nil
// 	case DRAM:
// 		return globalMSR.DRAM, nil
// 	case PP0:
// 		return globalMSR.PP0, nil
// 	case PP1:
// 		return globalMSR.PP1, nil
// 	}

// 	return domainMSRs{}, fmt.Errorf("No MSR for %s available", domain)
// }

// struct defs

//RAPLPowerLimitCtl specifies a power limit for a given time window
type RAPLPowerLimitCtl struct {
	PowerLimit      uint64
	EnableLimit     bool
	ClampingLimit   bool
	TimeWindowLimit uint64
}

//RAPLPowerLimit contains the data in the MSR_[DOMAIN]_POWER_LIMIT MSR
//This MSR containers two power limits. From the SDM:
//"Two power lmits can be specified, corresponding to time windows of different sizes"
//"Each power limit provides independent clamping control that would permit the processor cores to go below OS-requested state to meet the power limits."
type RAPLPowerLimit struct {
	Limit1 RAPLPowerLimitCtl
	Limit2 RAPLPowerLimitCtl
	Lock   bool
}
